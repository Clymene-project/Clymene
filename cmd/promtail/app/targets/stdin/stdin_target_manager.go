package stdin

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery/targetgroup"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

// bufferSize is the size of the buffered reader
const bufferSize = 8096

// file is an interface allowing us to abstract a file.
type file interface {
	Stat() (os.FileInfo, error)
	io.Reader
}

var (
	// stdIn is os.Stdin but can be replaced for testing purpose.
	stdIn       file = os.Stdin
	hostName, _      = os.Hostname()
	// defaultStdInCfg is the default config for stdin target if none provided.
	defaultStdInCfg = scrapeconfig.Config{
		JobName: "stdin",
		ServiceDiscoveryConfig: scrapeconfig.ServiceDiscoveryConfig{
			StaticConfigs: []*targetgroup.Group{
				{Labels: model.LabelSet{"job": "stdin"}},
				{Labels: model.LabelSet{"hostname": model.LabelValue(hostName)}},
			},
		},
	}
)

type Shutdownable interface {
	Shutdown()
}

// nolint:revive
type StdinTargetManager struct {
	*readerTarget
	app Shutdownable
}

func NewStdinTargetManager(reg prometheus.Registerer, log *zap.Logger, app Shutdownable, client api.EntryHandler, configs []scrapeconfig.Config) (*StdinTargetManager, error) {
	reader, err := newReaderTarget(reg, log, stdIn, client, getStdinConfig(log, configs))
	if err != nil {
		return nil, err
	}
	stdinManager := &StdinTargetManager{
		readerTarget: reader,
		app:          app,
	}
	// when we're done flushing our stdin we can shutdown the app.
	go func() {
		<-reader.ctx.Done()
		app.Shutdown()
	}()
	return stdinManager, nil
}

func getStdinConfig(log *zap.Logger, configs []scrapeconfig.Config) scrapeconfig.Config {
	cfg := defaultStdInCfg
	// if we receive configs we use the first one.
	if len(configs) > 0 {
		if len(configs) > 1 {
			log.Warn(fmt.Sprintf("too many scrape configs, skipping %d configs.", len(configs)-1))
		}
		cfg = configs[0]
	}
	return cfg
}

func (t *StdinTargetManager) Ready() bool {
	return t.ctx.Err() == nil
}
func (t *StdinTargetManager) Stop()                                     { t.cancel() }
func (t *StdinTargetManager) ActiveTargets() map[string][]target.Target { return nil }
func (t *StdinTargetManager) AllTargets() map[string][]target.Target    { return nil }

type readerTarget struct {
	in     *bufio.Reader
	out    api.EntryHandler
	lbs    model.LabelSet
	logger *zap.Logger

	cancel context.CancelFunc
	ctx    context.Context
}

func newReaderTarget(reg prometheus.Registerer, logger *zap.Logger, in io.Reader, client api.EntryHandler, cfg scrapeconfig.Config) (*readerTarget, error) {
	pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "pipeline")), cfg.PipelineStages, &cfg.JobName, reg)
	if err != nil {
		return nil, err
	}
	lbs := model.LabelSet{}
	for _, static := range cfg.ServiceDiscoveryConfig.StaticConfigs {
		if static != nil && static.Labels != nil {
			lbs = lbs.Merge(static.Labels)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	t := &readerTarget{
		in:     bufio.NewReaderSize(in, bufferSize),
		out:    pipeline.Wrap(client),
		cancel: cancel,
		ctx:    ctx,
		lbs:    lbs,
		logger: logger.With(zap.String("component", "reader")),
	}
	go t.read()

	return t, nil
}

func (t *readerTarget) read() {
	defer t.cancel()
	defer t.out.Stop()

	entries := t.out.Chan()
	for {
		if t.ctx.Err() != nil {
			return
		}
		line, err := t.in.ReadString('\n')
		if err != nil && err != io.EOF {
			t.logger.Warn("error reading buffer", zap.Error(err))
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if err == io.EOF {
				return
			}
			continue
		}
		entries <- api.Entry{
			Labels: t.lbs.Clone(),
			Entry: logproto.Entry{
				Timestamp: time.Now(),
				Line:      line,
			},
		}
		if err == io.EOF {
			return
		}
	}
}
