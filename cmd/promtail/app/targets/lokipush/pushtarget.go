package lokipush

import (
	"bufio"
	"flag"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/Clymene-project/Clymene/model/relabel"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	"github.com/Clymene-project/Clymene/pkg/tenant"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/weaveworks/common/server"
)

type PushTarget struct {
	logger        *zap.Logger
	handler       api.EntryHandler
	config        *scrapeconfig.PushTargetConfig
	relabelConfig []*relabel.Config
	jobName       string
	server        *server.Server
	registerer    prometheus.Registerer
}

func NewPushTarget(logger *zap.Logger,
	handler api.EntryHandler,
	relabel []*relabel.Config,
	jobName string,
	config *scrapeconfig.PushTargetConfig,
	reg prometheus.Registerer,
) (*PushTarget, error) {

	pt := &PushTarget{
		logger:        logger,
		handler:       handler,
		relabelConfig: relabel,
		jobName:       jobName,
		config:        config,
		registerer:    reg,
	}

	// Bit of a chicken and egg problem trying to register the defaults and apply overrides from the loaded config.
	// First create an empty config and set defaults.
	defaults := server.Config{}
	defaults.RegisterFlags(flag.NewFlagSet("empty", flag.ContinueOnError))
	// Then apply any config values loaded as overrides to the defaults.
	if err := mergo.Merge(&defaults, config.Server, mergo.WithOverride); err != nil {
		logger.Error("failed to parse configs and override defaults when configuring push server", zap.Error(err))
	}
	// The merge won't overwrite with a zero value but in the case of ports 0 value
	// indicates the desire for a random port so reset these to zero if the incoming config val is 0
	if config.Server.HTTPListenPort == 0 {
		defaults.HTTPListenPort = 0
	}
	if config.Server.GRPCListenPort == 0 {
		defaults.GRPCListenPort = 0
	}
	// Set the config to the new combined config.
	config.Server = defaults

	err := pt.run()
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (t *PushTarget) run() error {
	t.logger.Info("starting push server", zap.String("job", t.jobName))
	// To prevent metric collisions because all metrics are going to be registered in the global Prometheus registry.
	t.config.Server.MetricsNamespace = "promtail_" + t.jobName

	// We don't want the /debug and /metrics endpoints running
	t.config.Server.RegisterInstrumentation = false

	//util_log.InitLogger(&t.config.Server, t.registerer)

	srv, err := server.New(t.config.Server)
	if err != nil {
		return err
	}

	t.server = srv
	t.server.HTTP.Path("/loki/api/v1/push").Methods("POST").Handler(http.HandlerFunc(t.handleLoki))
	t.server.HTTP.Path("/promtail/api/v1/raw").Methods("POST").Handler(http.HandlerFunc(t.handlePlaintext))

	go func() {
		err := srv.Run()
		if err != nil {
			t.logger.Error("Loki push server shutdown with error", zap.Error(err))
		}
	}()

	return nil
}

func (t *PushTarget) handleLoki(w http.ResponseWriter, r *http.Request) {
	//_ = util_log.WithContext(r.Context(), util_log.Logger)
	_, _ = tenant.TenantID(r.Context())
	//req, err := push.ParseRequest(logger, userID, r, nil)
	//if err != nil {
	//	t.logger.Warn("failed to parse incoming push request", zap.Error(err))
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	var lastErr error
	//for _, stream := range req.Streams {
	//	ls, err := promql_parser.ParseMetric(stream.Labels)
	//	if err != nil {
	//		lastErr = err
	//		continue
	//	}
	//	sort.Sort(ls)
	//
	//	lb := labels.NewBuilder(ls)
	//
	//	// Add configured labels
	//	for k, v := range t.config.Labels {
	//		lb.Set(string(k), string(v))
	//	}
	//
	//	// Apply relabeling
	//	processed := relabel.Process(lb.Labels(), t.relabelConfig...)
	//	if processed == nil || len(processed) == 0 {
	//		w.WriteHeader(http.StatusNoContent)
	//		return
	//	}
	//
	//	// Convert to model.LabelSet
	//	filtered := model.LabelSet{}
	//	for i := range processed {
	//		if strings.HasPrefix(processed[i].Name, "__") {
	//			continue
	//		}
	//		filtered[model.LabelName(processed[i].Name)] = model.LabelValue(processed[i].Value)
	//	}
	//
	//	for _, entry := range stream.Entries {
	//		e := api.Entry{
	//			Labels: filtered.Clone(),
	//			Entry: logproto.Entry{
	//				Line: entry.Line,
	//			},
	//		}
	//		if t.config.KeepTimestamp {
	//			e.Timestamp = entry.Timestamp
	//		} else {
	//			e.Timestamp = time.Now()
	//		}
	//		t.handler.Chan() <- e
	//	}
	//}

	if lastErr != nil {
		t.logger.Warn("at least one entry in the push request failed to process", zap.Error(lastErr))
		http.Error(w, lastErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePlaintext handles newline delimited input such as plaintext or NDJSON.
func (t *PushTarget) handlePlaintext(w http.ResponseWriter, r *http.Request) {
	entries := t.handler.Chan()
	defer r.Body.Close()
	body := bufio.NewReader(r.Body)
	for {
		line, err := body.ReadString('\n')
		if err != nil && err != io.EOF {
			t.logger.Warn("failed to read incoming push request", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if err == io.EOF {
				break
			}
			continue
		}
		entries <- api.Entry{
			Labels: t.Labels().Clone(),
			Entry: logproto.Entry{
				Timestamp: time.Now(),
				Line:      line,
			},
		}
		if err == io.EOF {
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Type returns PushTargetType.
func (t *PushTarget) Type() target.TargetType {
	return target.PushTargetType
}

// Ready indicates whether or not the PushTarget target is ready to be read from.
func (t *PushTarget) Ready() bool {
	return true
}

// DiscoveredLabels returns the set of labels discovered by the PushTarget, which
// is always nil. Implements Target.
func (t *PushTarget) DiscoveredLabels() model.LabelSet {
	return nil
}

// Labels returns the set of labels that statically apply to all log entries
// produced by the PushTarget.
func (t *PushTarget) Labels() model.LabelSet {
	return t.config.Labels
}

// Details returns target-specific details.
func (t *PushTarget) Details() interface{} {
	return map[string]string{}
}

// Stop shuts down the PushTarget.
func (t *PushTarget) Stop() error {
	t.logger.Info("stopping push server", zap.String("job", t.jobName))
	t.server.Shutdown()
	t.handler.Stop()
	return nil
}
