//go:build cgo
// +build cgo

package journal

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
)

// JournalTargetManager manages a series of JournalTargets.
// nolint
type JournalTargetManager struct {
	logger  *zap.Logger
	targets map[string]*JournalTarget
}

// NewJournalTargetManager creates a new JournalTargetManager.
func NewJournalTargetManager(
	reg prometheus.Registerer,
	logger *zap.Logger,
	positions positions.Positions,
	client api.EntryHandler,
	scrapeConfigs []scrapeconfig.Config,
) (*JournalTargetManager, error) {
	tm := &JournalTargetManager{
		logger:  logger,
		targets: make(map[string]*JournalTarget),
	}

	for _, cfg := range scrapeConfigs {
		if cfg.JournalConfig == nil {
			continue
		}

		pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "journal_pipeline")), cfg.PipelineStages, &cfg.JobName, reg)
		if err != nil {
			return nil, err
		}

		t, err := NewJournalTarget(
			logger,
			pipeline.Wrap(client),
			positions,
			cfg.JobName,
			cfg.RelabelConfigs,
			cfg.JournalConfig,
		)
		if err != nil {
			return nil, err
		}

		tm.targets[cfg.JobName] = t
	}

	return tm, nil
}

// Ready returns true if at least one JournalTarget is also ready.
func (tm *JournalTargetManager) Ready() bool {
	for _, t := range tm.targets {
		if t.Ready() {
			return true
		}
	}
	return false
}

// Stop stops the JournalTargetManager and all of its JournalTargets.
func (tm *JournalTargetManager) Stop() {
	for _, t := range tm.targets {
		if err := t.Stop(); err != nil {
			t.logger.Error("error stopping JournalTarget", zap.Error(err))
		}
	}
}

// ActiveTargets returns the list of JournalTargets where journal data
// is being read. ActiveTargets is an alias to AllTargets as
// JournalTargets cannot be deactivated, only stopped.
func (tm *JournalTargetManager) ActiveTargets() map[string][]target.Target {
	return tm.AllTargets()
}

// AllTargets returns the list of all targets where journal data
// is currently being read.
func (tm *JournalTargetManager) AllTargets() map[string][]target.Target {
	result := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		result[k] = []target.Target{v}
	}
	return result
}
