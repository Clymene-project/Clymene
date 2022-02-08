package docker

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery/targetgroup"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/Clymene-project/Clymene/model/relabel"
	"go.uber.org/zap"
	"sync"

	"github.com/docker/docker/client"
	"github.com/prometheus/common/model"
)

const DockerSource = "Docker"

// targetGroup manages all container targets of one Docker daemon.
type targetGroup struct {
	metrics       *Metrics
	logger        *zap.Logger
	positions     positions.Positions
	entryHandler  api.EntryHandler
	defaultLabels model.LabelSet
	relabelConfig []*relabel.Config
	host          string
	client        client.APIClient

	mtx     sync.Mutex
	targets map[string]*Target
}

func (tg *targetGroup) sync(groups []*targetgroup.Group) {
	tg.mtx.Lock()
	defer tg.mtx.Unlock()

	for _, group := range groups {
		if group.Source != DockerSource {
			continue
		}

		for _, t := range group.Targets {
			containerID, ok := t[dockerLabelContainerID]
			if !ok {
				tg.logger.Debug("Docker target did not include container ID")
				continue
			}

			err := tg.addTarget(string(containerID), t)
			if err != nil {
				tg.logger.Error("could not add target", zap.String("containerID", string(containerID)), zap.Error(err))
			}
		}
	}
}

// addTarget checks whether the container with given id is already known. If not it's added to the this group
func (tg *targetGroup) addTarget(id string, discoveredLabels model.LabelSet) error {
	if tg.client == nil {
		var err error
		opts := []client.Opt{
			client.WithHost(tg.host),
			client.WithAPIVersionNegotiation(),
		}
		tg.client, err = client.NewClientWithOpts(opts...)
		if err != nil {
			tg.logger.Error("could not create new Docker client", zap.Error(err))
			return err
		}
	}

	if t, ok := tg.targets[id]; ok {
		tg.logger.Debug("container target already exists", zap.String("container", id))
		t.startIfNotRunning()
		return nil
	}

	t, err := NewTarget(
		tg.metrics,
		tg.logger.With(zap.String("target", fmt.Sprintf("docker/%s", id))),
		tg.entryHandler,
		tg.positions,
		id,
		discoveredLabels.Merge(tg.defaultLabels),
		tg.relabelConfig,
		tg.client,
	)
	if err != nil {
		return err
	}
	tg.targets[id] = t
	tg.logger.Error("added Docker target", zap.String("container", id))
	return nil
}

// Ready returns true if at least one target is running.
func (tg *targetGroup) Ready() bool {
	tg.mtx.Lock()
	defer tg.mtx.Unlock()

	for _, t := range tg.targets {
		if t.Ready() {
			return true
		}
	}

	return true
}

// Stop all targets
func (tg *targetGroup) Stop() {
	tg.mtx.Lock()
	defer tg.mtx.Unlock()

	for _, t := range tg.targets {
		t.Stop()
	}
	tg.entryHandler.Stop()
}

// ActiveTargets return all targets that are ready.
func (tg *targetGroup) ActiveTargets() []target.Target {
	tg.mtx.Lock()
	defer tg.mtx.Unlock()

	result := make([]target.Target, 0, len(tg.targets))
	for _, t := range tg.targets {
		if t.Ready() {
			result = append(result, t)
		}
	}
	return result
}

// AllTargets returns all targets of this group.
func (tg *targetGroup) AllTargets() []target.Target {
	result := make([]target.Target, 0, len(tg.targets))
	for _, t := range tg.targets {
		result = append(result, t)
	}
	return result
}
