package stages

import (
	"fmt"
	"github.com/go-logfmt/logfmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

// Config Errors
const (
	ErrMappingRequired        = "logfmt mapping is required"
	ErrEmptyLogfmtStageConfig = "empty logfmt stage configuration"
	ErrEmptyLogfmtStageSource = "empty source"
)

// LogfmtConfig represents a logfmt Stage configuration
type LogfmtConfig struct {
	Mapping map[string]string `mapstructure:"mapping"`
	Source  *string           `mapstructure:"source"`
}

// validateLogfmtConfig validates a logfmt stage config and returns an inverse mapping of configured mapping.
// Mapping inverse is done to make lookup easier. The key would be the key from parsed logfmt and
// value would be the key with which the data in extracted map would be set.
func validateLogfmtConfig(c *LogfmtConfig) (map[string]string, error) {
	if c == nil {
		return nil, errors.New(ErrEmptyLogfmtStageConfig)
	}

	if len(c.Mapping) == 0 {
		return nil, errors.New(ErrMappingRequired)
	}

	if c.Source != nil && *c.Source == "" {
		return nil, errors.New(ErrEmptyLogfmtStageSource)
	}

	inverseMapping := make(map[string]string)
	for k, v := range c.Mapping {
		// if value is not set, use the key for setting data in extracted map.
		if v == "" {
			v = k
		}
		inverseMapping[v] = k
	}

	return inverseMapping, nil
}

// logfmtStage sets extracted data using logfmt parser
type logfmtStage struct {
	cfg            *LogfmtConfig
	inverseMapping map[string]string
	logger         *zap.Logger
}

// newLogfmtStage creates a new logfmt pipeline stage from a config.
func newLogfmtStage(logger *zap.Logger, config interface{}) (Stage, error) {
	cfg, err := parseLogfmtConfig(config)
	if err != nil {
		return nil, err
	}

	// inverseMapping would hold the mapping in inverse which would make lookup easier.
	// To explain it simply, the key would be the key from parsed logfmt and value would be the key with which the data in extracted map would be set.
	inverseMapping, err := validateLogfmtConfig(cfg)
	if err != nil {
		return nil, err
	}

	return toStage(&logfmtStage{
		cfg:            cfg,
		inverseMapping: inverseMapping,
		logger:         logger.With(zap.String("component", "stage"), zap.String("type", "logfmt")),
	}), nil
}

func parseLogfmtConfig(config interface{}) (*LogfmtConfig, error) {
	cfg := &LogfmtConfig{}
	err := mapstructure.Decode(config, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Process implements Stage
func (j *logfmtStage) Process(labels model.LabelSet, extracted map[string]interface{}, t *time.Time, entry *string) {
	// If a source key is provided, the logfmt stage should process it
	// from the extracted map, otherwise should fallback to the entry
	input := entry

	if j.cfg.Source != nil {
		if _, ok := extracted[*j.cfg.Source]; !ok {
			if Debug {
				j.logger.Debug("source does not exist in the set of extracted values", zap.String("source", *j.cfg.Source))
			}
			return
		}

		value, err := getString(extracted[*j.cfg.Source])
		if err != nil {
			if Debug {
				j.logger.Debug("failed to convert source value to string", zap.String("source", *j.cfg.Source), zap.Error(err), zap.String("type", reflect.TypeOf(extracted[*j.cfg.Source]).String()))
			}
			return
		}

		input = &value
	}

	if input == nil {
		if Debug {
			j.logger.Debug("cannot parse a nil entry")
		}
		return
	}
	decoder := logfmt.NewDecoder(strings.NewReader(*input))
	extractedEntriesCount := 0
	for decoder.ScanRecord() {
		for decoder.ScanKeyval() {
			mapKey, ok := j.inverseMapping[string(decoder.Key())]
			if ok {
				extracted[mapKey] = string(decoder.Value())
				extractedEntriesCount++
			}
		}
	}

	if decoder.Err() != nil {
		j.logger.Error("failed to decode logfmt", zap.Error(decoder.Err()))
		return
	}

	if Debug {
		if extractedEntriesCount != len(j.inverseMapping) {
			j.logger.Debug(fmt.Sprintf("found only %d out of %d configured mappings in logfmt stage", extractedEntriesCount, len(j.inverseMapping)))
		}
		j.logger.Debug("extracted data debug in logfmt stage", zap.String("extracted data", fmt.Sprintf("%v", extracted)))
	}
}

// Name implements Stage
func (j *logfmtStage) Name() string {
	return StageTypeLogfmt
}
