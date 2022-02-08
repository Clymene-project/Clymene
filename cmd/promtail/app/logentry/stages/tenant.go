package stages

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"go.uber.org/zap"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

const (
	ErrTenantStageEmptySourceOrValue        = "source or value config are required"
	ErrTenantStageConflictingSourceAndValue = "source and value are mutually exclusive: you should set source or value but not both"
)

type tenantStage struct {
	cfg    TenantConfig
	logger *zap.Logger
}

type TenantConfig struct {
	Source string `mapstructure:"source"`
	Value  string `mapstructure:"value"`
}

// validateTenantConfig validates the tenant stage configuration
func validateTenantConfig(c TenantConfig) error {
	if c.Source == "" && c.Value == "" {
		return errors.New(ErrTenantStageEmptySourceOrValue)
	}

	if c.Source != "" && c.Value != "" {
		return errors.New(ErrTenantStageConflictingSourceAndValue)
	}

	return nil
}

// newTenantStage creates a new tenant stage to override the tenant ID from extracted data
func newTenantStage(logger *zap.Logger, configs interface{}) (Stage, error) {
	cfg := TenantConfig{}
	err := mapstructure.Decode(configs, &cfg)
	if err != nil {
		return nil, err
	}

	err = validateTenantConfig(cfg)
	if err != nil {
		return nil, err
	}

	return toStage(&tenantStage{
		cfg:    cfg,
		logger: logger,
	}), nil
}

// Process implements Stage
func (s *tenantStage) Process(labels model.LabelSet, extracted map[string]interface{}, t *time.Time, entry *string) {
	var tenantID string

	// Get tenant ID from source or configured value
	if s.cfg.Source != "" {
		tenantID = s.getTenantFromSourceField(extracted)
	} else {
		tenantID = s.cfg.Value
	}

	// Skip an empty tenant ID (ie. failed to get the tenant from the source)
	if tenantID == "" {
		return
	}

	labels[client.ReservedLabelTenantID] = model.LabelValue(tenantID)
}

// Name implements Stage
func (s *tenantStage) Name() string {
	return StageTypeTenant
}

func (s *tenantStage) getTenantFromSourceField(extracted map[string]interface{}) string {
	// Get the tenant ID from the source data
	value, ok := extracted[s.cfg.Source]
	if !ok {
		if Debug {
			s.logger.Debug("the tenant source does not exist in the extracted data", zap.String("source", s.cfg.Source))
		}
		return ""
	}

	// Convert the value to string
	tenantID, err := getString(value)
	if err != nil {
		if Debug {
			s.logger.Debug("failed to convert value to string", zap.Error(err), zap.String("type", reflect.TypeOf(value).String()))
		}
		return ""
	}

	return tenantID
}
