package cloudflare

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"go.uber.org/zap"
)

// TargetManager manages a series of cloudflare targets.
type TargetManager struct {
	logger  *zap.Logger
	targets map[string]*Target
}

// NewTargetManager creates a new cloudflare target managers.
func NewTargetManager(
	metrics *Metrics,
	logger *zap.Logger,
	positions positions.Positions,
	pushClient api.EntryHandler,
	scrapeConfigs []scrapeconfig.Config,
) (*TargetManager, error) {
	tm := &TargetManager{
		logger:  logger,
		targets: make(map[string]*Target),
	}
	for _, cfg := range scrapeConfigs {
		if cfg.CloudflareConfig == nil {
			continue
		}
		pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "cloudflare_pipeline")), cfg.PipelineStages, &cfg.JobName, metrics.reg)
		if err != nil {
			return nil, err
		}
		t, err := NewTarget(metrics, logger.With(zap.String("target", "cloudflare")), pipeline.Wrap(pushClient), positions, cfg.CloudflareConfig)
		if err != nil {
			return nil, err
		}
		tm.targets[cfg.JobName] = t
	}

	return tm, nil
}

// Ready returns true if at least one cloudflare target is active.
func (tm *TargetManager) Ready() bool {
	for _, t := range tm.targets {
		if t.Ready() {
			return true
		}
	}
	return false
}

func (tm *TargetManager) Stop() {
	for _, t := range tm.targets {
		t.Stop()
	}
}

func (tm *TargetManager) ActiveTargets() map[string][]target.Target {
	result := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		if v.Ready() {
			result[k] = []target.Target{v}
		}
	}
	return result
}

func (tm *TargetManager) AllTargets() map[string][]target.Target {
	result := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		result[k] = []target.Target{v}
	}
	return result
}
