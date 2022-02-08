package stages

import (
	"bytes"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	"go.uber.org/zap"
	"regexp"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

const (
	ErrMultilineStageEmptyConfig        = "multiline stage config must define `firstline` regular expression"
	ErrMultilineStageInvalidRegex       = "multiline stage first line regex compilation error: %v"
	ErrMultilineStageInvalidMaxWaitTime = "multiline stage `max_wait_time` parse error: %v"
)

const (
	maxLineDefault uint64 = 128
	maxWaitDefault        = 3 * time.Second
)

// MultilineConfig contains the configuration for a multilineStage
type MultilineConfig struct {
	Expression  *string `mapstructure:"firstline"`
	regex       *regexp.Regexp
	MaxLines    *uint64 `mapstructure:"max_lines"`
	MaxWaitTime *string `mapstructure:"max_wait_time"`
	maxWait     time.Duration
}

func validateMultilineConfig(cfg *MultilineConfig) error {
	if cfg == nil || cfg.Expression == nil {
		return errors.New(ErrMultilineStageEmptyConfig)
	}

	expr, err := regexp.Compile(*cfg.Expression)
	if err != nil {
		return errors.Errorf(ErrMultilineStageInvalidRegex, err)
	}
	cfg.regex = expr

	if cfg.MaxWaitTime != nil {
		maxWait, err := time.ParseDuration(*cfg.MaxWaitTime)
		if err != nil {
			return errors.Errorf(ErrMultilineStageInvalidMaxWaitTime, err)
		}
		cfg.maxWait = maxWait
	} else {
		cfg.maxWait = maxWaitDefault
	}

	if cfg.MaxLines == nil {
		cfg.MaxLines = new(uint64)
		*cfg.MaxLines = maxLineDefault
	}

	return nil
}

// multilineStage matches lines to determine whether the following lines belong to a block and should be collapsed
type multilineStage struct {
	logger *zap.Logger
	cfg    *MultilineConfig
}

// multilineState captures the internal state of a running multiline stage.
type multilineState struct {
	buffer         *bytes.Buffer // The lines of the current multiline block.
	startLineEntry Entry         // The entry of the start line of a multiline block.
	currentLines   uint64        // The number of lines of the current multiline block.
}

// newMulitlineStage creates a MulitlineStage from config
func newMultilineStage(logger *zap.Logger, config interface{}) (Stage, error) {
	cfg := &MultilineConfig{}
	err := mapstructure.WeakDecode(config, cfg)
	if err != nil {
		return nil, err
	}
	err = validateMultilineConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &multilineStage{
		logger: logger.With(zap.String("component", "stage"), zap.String("type", "multiline")),
		cfg:    cfg,
	}, nil
}

func (m *multilineStage) Run(in chan Entry) chan Entry {
	out := make(chan Entry)
	go func() {
		defer close(out)

		streams := make(map[model.Fingerprint](chan Entry))
		wg := new(sync.WaitGroup)

		for e := range in {
			key := e.Labels.FastFingerprint()
			s, ok := streams[key]
			if !ok {
				// Pass through entries until we hit first start line.
				if !m.cfg.regex.MatchString(e.Line) {
					if Debug {
						m.logger.Debug("pass through entry", zap.String("stream", key.String()))
					}
					out <- e
					continue
				}

				if Debug {
					m.logger.Debug("creating new stream", zap.String("stream", key.String()))
				}
				s = make(chan Entry)
				streams[key] = s

				wg.Add(1)
				go m.runMultiline(s, out, wg)
			}
			if Debug {
				m.logger.Debug("pass entry", zap.String("stream", key.String()), zap.String("line", e.Line))
			}
			s <- e
		}

		// Close all streams and wait for them to finish being processed.
		for _, s := range streams {
			close(s)
		}
		wg.Wait()
	}()
	return out
}

func (m *multilineStage) runMultiline(in chan Entry, out chan Entry, wg *sync.WaitGroup) {
	defer wg.Done()

	state := &multilineState{
		buffer:       new(bytes.Buffer),
		currentLines: 0,
	}

	for {
		select {
		case <-time.After(m.cfg.maxWait):
			if Debug {
				m.logger.Debug(fmt.Sprintf("flush multiline block due to %v timeout", m.cfg.maxWait), zap.String("block", state.buffer.String()))
			}
			m.flush(out, state)
		case e, ok := <-in:
			if Debug {
				m.logger.Debug("processing line", zap.String("line", e.Line), zap.String("stream", e.Labels.FastFingerprint().String()))
			}

			if !ok {
				if Debug {
					m.logger.Debug("flush multiline block because inbound closed", zap.String("block", state.buffer.String()), zap.String("stream", e.Labels.FastFingerprint().String()))
				}
				m.flush(out, state)
				return
			}

			isFirstLine := m.cfg.regex.MatchString(e.Line)
			if isFirstLine {
				if Debug {
					m.logger.Debug("flush multiline block because new start line", zap.String("block", state.buffer.String()), zap.String("stream", e.Labels.FastFingerprint().String()))
				}
				m.flush(out, state)

				// The start line entry is used to set timestamp and labels in the flush method.
				// The timestamps for following lines are ignored for now.
				state.startLineEntry = e
			}

			// Append block line
			if state.buffer.Len() > 0 {
				state.buffer.WriteRune('\n')
			}
			state.buffer.WriteString(e.Line)
			state.currentLines++

			if state.currentLines == *m.cfg.MaxLines {
				m.flush(out, state)
			}
		}
	}
}

func (m *multilineStage) flush(out chan Entry, s *multilineState) {
	if s.buffer.Len() == 0 {
		if Debug {
			m.logger.Debug("nothing to flush", zap.Int("buffer_len", s.buffer.Len()))
		}
		return
	}
	// copy extracted data.
	extracted := make(map[string]interface{}, len(s.startLineEntry.Extracted))
	for k, v := range s.startLineEntry.Extracted {
		extracted[k] = v
	}
	collapsed := Entry{
		Extracted: extracted,
		Entry: api.Entry{
			Labels: s.startLineEntry.Entry.Labels.Clone(),
			Entry: logproto.Entry{
				Timestamp: s.startLineEntry.Entry.Entry.Timestamp,
				Line:      s.buffer.String(),
			},
		},
	}
	s.buffer.Reset()
	s.currentLines = 0

	out <- collapsed
}

// Name implements Stage
func (m *multilineStage) Name() string {
	return StageTypeMultiline
}
