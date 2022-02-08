//go:build windows
// +build windows

package windows

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// TargetManager manages a series of windows event targets.
type TargetManager struct {
	logger  *zap.Logger
	targets map[string]*Target
}

// NewTargetManager creates a new Windows managers.
func NewTargetManager(
	reg prometheus.Registerer,
	logger *zap.Logger,
	client api.EntryHandler,
	scrapeConfigs []scrapeconfig.Config,
) (*TargetManager, error) {
	tm := &TargetManager{
		logger:  logger,
		targets: make(map[string]*Target),
	}

	for _, cfg := range scrapeConfigs {
		pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "windows_pipeline")), cfg.PipelineStages, &cfg.JobName, reg)
		if err != nil {
			return nil, err
		}

		t, err := New(logger, pipeline.Wrap(client), cfg.RelabelConfigs, cfg.WindowsConfig)
		if err != nil {
			return nil, err
		}

		tm.targets[cfg.JobName] = t
	}

	return tm, nil
}

// Ready returns true if at least one Windows target is also ready.
func (tm *TargetManager) Ready() bool {
	for _, t := range tm.targets {
		if t.Ready() {
			return true
		}
	}
	return false
}

// Stop stops the Windows target manager and all of its targets.
func (tm *TargetManager) Stop() {
	for _, t := range tm.targets {
		if err := t.Stop(); err != nil {
			t.logger.Error("error stopping windows target", zap.Error(err))
		}
	}
}

// ActiveTargets returns the list of active Windows targets.
func (tm *TargetManager) ActiveTargets() map[string][]target.Target {
	result := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		if v.Ready() {
			result[k] = []target.Target{v}
		}
	}
	return result
}

// AllTargets returns the list of all targets.
func (tm *TargetManager) AllTargets() map[string][]target.Target {
	result := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		result[k] = []target.Target{v}
	}
	return result
}
