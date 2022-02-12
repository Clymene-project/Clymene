package file

import (
	"context"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery/kubernetes"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery/targetgroup"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/Clymene-project/Clymene/model/labels"
	"github.com/Clymene-project/Clymene/model/relabel"
	util "github.com/Clymene-project/Clymene/pkg/lokiutil"
	"go.uber.org/zap"
	"gopkg.in/fsnotify.v1"
	"os"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

const (
	pathLabel              = "__path__"
	hostLabel              = "__host__"
	kubernetesPodNodeField = "spec.nodeName"
)

// FileTargetManager manages a set of targets.
// nolint:revive
type FileTargetManager struct {
	log     *zap.Logger
	quit    context.CancelFunc
	syncers map[string]*targetSyncer
	manager *discovery.Manager

	watcher            *fsnotify.Watcher
	targetEventHandler chan fileTargetEvent

	wg sync.WaitGroup
}

// NewFileTargetManager creates a new TargetManager.
func NewFileTargetManager(
	metrics *Metrics,
	logger *zap.Logger,
	positions positions.Positions,
	client api.EntryHandler,
	scrapeConfigs []scrapeconfig.Config,
	targetConfig *Config,
) (*FileTargetManager, error) {
	reg := metrics.reg
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	ctx, quit := context.WithCancel(context.Background())
	tm := &FileTargetManager{
		log:                logger,
		quit:               quit,
		watcher:            watcher,
		targetEventHandler: make(chan fileTargetEvent),
		syncers:            map[string]*targetSyncer{},
		manager:            discovery.NewManager(ctx, logger.With(zap.String("component", "discovery"))),
	}

	hostname, err := hostname()
	if err != nil {
		return nil, err
	}

	configs := map[string]discovery.Configs{}
	for _, cfg := range scrapeConfigs {
		if !cfg.HasServiceDiscoveryConfig() {
			continue
		}

		pipeline, err := stages.NewPipeline(logger.With(zap.String("component", "file_pipeline")), cfg.PipelineStages, &cfg.JobName, reg)
		if err != nil {
			return nil, err
		}

		// Add Source value to the static config target groups for unique identification
		// within scrape pool. Also, default target label to localhost if target is not
		// defined in promtail config.
		// Just to make sure prometheus target group sync works fine.
		for i, tg := range cfg.ServiceDiscoveryConfig.StaticConfigs {
			tg.Source = fmt.Sprintf("%d", i)
			if len(tg.Targets) == 0 {
				tg.Targets = []model.LabelSet{
					{model.AddressLabel: "localhost"},
				}
			}
		}

		// Add an additional api-level node filtering, so we only fetch pod metadata for
		// all the pods from the current node. Without this filtering we will have to
		// download metadata for all pods running on a cluster, which may be a long operation.
		for _, kube := range cfg.ServiceDiscoveryConfig.KubernetesSDConfigs {
			if kube.Role == kubernetes.RolePod {
				selector := fmt.Sprintf("%s=%s", kubernetesPodNodeField, hostname)
				kube.Selectors = []kubernetes.SelectorConfig{
					{Role: kubernetes.RolePod, Field: selector},
				}
			}
		}

		s := &targetSyncer{
			metrics:           metrics,
			log:               logger,
			positions:         positions,
			relabelConfig:     cfg.RelabelConfigs,
			targets:           map[string]*FileTarget{},
			droppedTargets:    []target.Target{},
			hostname:          hostname,
			entryHandler:      pipeline.Wrap(client),
			targetConfig:      targetConfig,
			fileEventWatchers: map[string]chan fsnotify.Event{},
		}
		tm.syncers[cfg.JobName] = s
		configs[cfg.JobName] = cfg.ServiceDiscoveryConfig.Configs()
	}

	tm.wg.Add(3)
	go tm.run(ctx)
	go tm.watchTargetEvents(ctx)
	go tm.watchFsEvents(ctx)

	go util.LogError("running target manager", logger, tm.manager.Run)

	return tm, tm.manager.ApplyConfig(configs)
}

func (tm *FileTargetManager) watchTargetEvents(ctx context.Context) {
	defer tm.wg.Done()

	for {
		select {
		case event := <-tm.targetEventHandler:
			switch event.eventType {
			case fileTargetEventWatchStart:
				if err := tm.watcher.Add(event.path); err != nil {
					tm.log.Error("error adding directory to watcher", zap.Error(err))
				}
			case fileTargetEventWatchStop:
				if err := tm.watcher.Remove(event.path); err != nil {
					tm.log.Error(" failed to remove directory from watcher", zap.Error(err))
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (tm *FileTargetManager) watchFsEvents(ctx context.Context) {
	defer tm.wg.Done()

	for {
		select {
		case event := <-tm.watcher.Events:
			// we only care about Create events
			if event.Op == fsnotify.Create {
				tm.log.Info("received file watcher event", zap.String("name", event.Name), zap.String("op", event.Op.String()))
				for _, s := range tm.syncers {
					s.sendFileCreateEvent(event)
				}
			}
		case err := <-tm.watcher.Errors:
			tm.log.Error("error from fswatch", zap.Error(err))
		case <-ctx.Done():
			return
		}
	}
}

func (tm *FileTargetManager) run(ctx context.Context) {
	defer tm.wg.Done()

	for {
		select {
		case targetGroups := <-tm.manager.SyncCh():
			for jobName, groups := range targetGroups {
				tm.syncers[jobName].sync(groups, tm.targetEventHandler)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Ready if there's at least one file target
func (tm *FileTargetManager) Ready() bool {
	for _, s := range tm.syncers {
		if s.ready() {
			return true
		}
	}
	return false
}

// Stop the TargetManager.
func (tm *FileTargetManager) Stop() {
	tm.quit()
	tm.wg.Wait()

	for _, s := range tm.syncers {
		s.stop()
	}
	util.LogError("closing watcher", tm.log, tm.watcher.Close)
	close(tm.targetEventHandler)
}

// ActiveTargets returns the active targets currently being scraped.
func (tm *FileTargetManager) ActiveTargets() map[string][]target.Target {
	result := map[string][]target.Target{}
	for jobName, syncer := range tm.syncers {
		result[jobName] = append(result[jobName], syncer.ActiveTargets()...)
	}
	return result
}

// AllTargets returns all targets, active and dropped.
func (tm *FileTargetManager) AllTargets() map[string][]target.Target {
	result := map[string][]target.Target{}
	for jobName, syncer := range tm.syncers {
		result[jobName] = append(result[jobName], syncer.ActiveTargets()...)
		result[jobName] = append(result[jobName], syncer.DroppedTargets()...)
	}
	return result
}

// targetSyncer sync targets based on service discovery changes.
type targetSyncer struct {
	metrics      *Metrics
	log          *zap.Logger
	positions    positions.Positions
	entryHandler api.EntryHandler
	hostname     string

	fileEventWatchers map[string]chan fsnotify.Event

	droppedTargets []target.Target
	targets        map[string]*FileTarget
	mtx            sync.Mutex

	relabelConfig []*relabel.Config
	targetConfig  *Config
}

// sync synchronize target based on received target groups received by service discovery
func (s *targetSyncer) sync(groups []*targetgroup.Group, targetEventHandler chan fileTargetEvent) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	targets := map[string]struct{}{}
	dropped := []target.Target{}

	for _, group := range groups {
		for _, t := range group.Targets {
			s.log.Debug("new target", zap.String("labels", t.String()))

			discoveredLabels := group.Labels.Merge(t)
			var labelMap = make(map[string]string)
			for k, v := range discoveredLabels.Clone() {
				labelMap[string(k)] = string(v)
			}

			processedLabels := relabel.Process(labels.FromMap(labelMap), s.relabelConfig...)

			var labels = make(model.LabelSet)
			for k, v := range processedLabels.Map() {
				labels[model.LabelName(k)] = model.LabelValue(v)
			}

			// Drop empty targets (drop in relabeling).
			if processedLabels == nil {
				dropped = append(dropped, target.NewDroppedTarget("dropping target, no labels", discoveredLabels))
				s.log.Debug("dropping target, no labels")
				s.metrics.failedTargets.WithLabelValues("empty_labels").Inc()
				continue
			}

			host, ok := labels[hostLabel]
			if ok && string(host) != s.hostname {
				dropped = append(dropped, target.NewDroppedTarget(fmt.Sprintf("ignoring target, wrong host (labels:%s hostname:%s)", labels.String(), s.hostname), discoveredLabels))
				s.log.Debug("ignoring target, wrong host", zap.String("labels", labels.String()), zap.String("hostname", s.hostname))
				s.metrics.failedTargets.WithLabelValues("wrong_host").Inc()
				continue
			}

			path, ok := labels[pathLabel]
			if !ok {
				dropped = append(dropped, target.NewDroppedTarget("no path for target", discoveredLabels))
				s.log.Info("no path for target", zap.String("labels", labels.String()))
				s.metrics.failedTargets.WithLabelValues("no_path").Inc()
				continue
			}

			for k := range labels {
				if strings.HasPrefix(string(k), "__") {
					delete(labels, k)
				}
			}

			key := fmt.Sprintf("%s:%s", path, labels.String())
			targets[key] = struct{}{}
			if _, ok := s.targets[key]; ok {
				dropped = append(dropped, target.NewDroppedTarget("ignoring target, already exists", discoveredLabels))
				s.log.Debug("ignoring target, already exists", zap.String("labels", labels.String()))
				s.metrics.failedTargets.WithLabelValues("exists").Inc()
				continue
			}

			s.log.Info("Adding target", zap.String("key", key))

			wkey := string(path)
			watcher, ok := s.fileEventWatchers[wkey]
			if !ok {
				watcher = make(chan fsnotify.Event)
				s.fileEventWatchers[wkey] = watcher
			}
			t, err := s.newTarget(wkey, labels, discoveredLabels, watcher, targetEventHandler)
			if err != nil {
				dropped = append(dropped, target.NewDroppedTarget(fmt.Sprintf("Failed to create target: %s", err.Error()), discoveredLabels))
				s.log.Error("Failed to create target", zap.String("key", key), zap.Error(err))
				s.metrics.failedTargets.WithLabelValues("error").Inc()
				continue
			}

			s.metrics.targetsActive.Add(1.)
			s.targets[key] = t
		}
	}

	for key, target := range s.targets {
		if _, ok := targets[key]; !ok {
			s.log.Info("Removing target", zap.String("key", key))
			target.Stop()
			s.metrics.targetsActive.Add(-1.)
			delete(s.targets, key)

			// close related file event watcher
			k := target.path
			if _, ok := s.fileEventWatchers[k]; ok {
				close(s.fileEventWatchers[k])
				delete(s.fileEventWatchers, k)
			} else {
				s.log.Warn("failed to remove file event watcher", zap.String("path", k))
			}
		}
	}
	s.droppedTargets = dropped
}

// sendFileCreateEvent sends file creation events to only the targets with matched path.
func (s *targetSyncer) sendFileCreateEvent(event fsnotify.Event) {
	// Lock the mutex because other threads are manipulating s.fileEventWatchers which can lead to a deadlock
	// where we send events to channels where nobody is listening anymore
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for path, watcher := range s.fileEventWatchers {
		matched, err := doublestar.Match(path, event.Name)
		if err != nil {
			s.log.Error("failed to match file", zap.Error(err), zap.String("filename", event.Name))
			continue
		}
		if !matched {
			s.log.Debug("new file does not match glob", zap.String("filename", event.Name))
			continue
		}
		watcher <- event
	}
}

func (s *targetSyncer) newTarget(path string, labels model.LabelSet, discoveredLabels model.LabelSet, fileEventWatcher chan fsnotify.Event, targetEventHandler chan fileTargetEvent) (*FileTarget, error) {
	return NewFileTarget(s.metrics, s.log, s.entryHandler, s.positions, path, labels, discoveredLabels, s.targetConfig, fileEventWatcher, targetEventHandler)
}

func (s *targetSyncer) DroppedTargets() []target.Target {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return append([]target.Target(nil), s.droppedTargets...)
}

func (s *targetSyncer) ActiveTargets() []target.Target {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	actives := []target.Target{}
	for _, t := range s.targets {
		actives = append(actives, t)
	}
	return actives
}

func (s *targetSyncer) ready() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, target := range s.targets {
		if target.Ready() {
			return true
		}
	}
	return false
}

func (s *targetSyncer) stop() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for key, target := range s.targets {
		s.log.Info("Removing target", zap.String("key", key))
		target.Stop()
		delete(s.targets, key)
	}

	for key, watcher := range s.fileEventWatchers {
		close(watcher)
		delete(s.fileEventWatchers, key)
	}
	s.entryHandler.Stop()
}

func hostname() (string, error) {
	hostname := os.Getenv("HOSTNAME")
	if hostname != "" {
		return hostname, nil
	}

	return os.Hostname()
}
