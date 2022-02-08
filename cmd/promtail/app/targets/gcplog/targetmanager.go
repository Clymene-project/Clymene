package gcplog

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"go.uber.org/zap"
)

// nolint:revive
type GcplogTargetManager struct {
	logger  *zap.Logger
	targets map[string]*GcplogTarget
}

func NewGcplogTargetManager(
	metrics *Metrics,
	logger *zap.Logger,
	client api.EntryHandler,
	scrape []scrapeconfig.Config,
) (*GcplogTargetManager, error) {
	tm := &GcplogTargetManager{
		logger:  logger,
		targets: make(map[string]*GcplogTarget),
	}

	for _, cf := range scrape {
		if cf.GcplogConfig == nil {
			continue
		}
		pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "pubsub_pipeline")), cf.PipelineStages, &cf.JobName, metrics.reg)
		if err != nil {
			return nil, err
		}

		t, err := NewGcplogTarget(metrics, logger, pipeline.Wrap(client), cf.RelabelConfigs, cf.JobName, cf.GcplogConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create pubsub target: %w", err)
		}
		tm.targets[cf.JobName] = t
	}

	return tm, nil
}

func (tm *GcplogTargetManager) Ready() bool {
	for _, t := range tm.targets {
		if t.Ready() {
			return true
		}
	}
	return false
}

func (tm *GcplogTargetManager) Stop() {
	for name, t := range tm.targets {
		if err := t.Stop(); err != nil {
			t.logger.Error("failed to stop pubsub target", zap.String("name", name), zap.Error(err))
		}
	}
}

func (tm *GcplogTargetManager) ActiveTargets() map[string][]target.Target {
	// TODO(kavi): if someway to check if specific topic is active and store the state on the target struct?
	return tm.AllTargets()
}

func (tm *GcplogTargetManager) AllTargets() map[string][]target.Target {
	res := make(map[string][]target.Target, len(tm.targets))
	for k, v := range tm.targets {
		res[k] = []target.Target{v}
	}
	return res
}
