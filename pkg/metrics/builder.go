package metrics

import (
	"errors"
	"expvar"
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

const (
	metricsBackend        = "metrics-backend"
	metricsHTTPRoute      = "metrics-http-route"
	defaultMetricsBackend = "prometheus"
	defaultMetricsRoute   = "/metrics"
)

var errUnknownBackend = errors.New("unknown metrics backend specified")

// Builder provides command line options to configure metrics backend used by Jaeger executables.
type Builder struct {
	Backend   string
	HTTPRoute string // endpoint name to expose metrics, e.g. for scraping
	handler   http.Handler
}

// AddFlags adds flags for Builder.
func AddFlags(flags *flag.FlagSet) {
	flags.String(
		metricsBackend,
		defaultMetricsBackend,
		"Defines which metrics backend to use for metrics reporting: expvar, prometheus, none")
	flags.String(
		metricsHTTPRoute,
		defaultMetricsRoute,
		"Defines the route of HTTP endpoint for metrics backends that support scraping")
}

// InitFromViper initializes Builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.Backend = v.GetString(metricsBackend)
	b.HTTPRoute = v.GetString(metricsHTTPRoute)
	return b
}

// CreateMetricsFactory creates a metrics factory based on the configured type of the backend.
// If the metrics backend supports HTTP endpoint for scraping, it is stored in the builder and
// can be later added by RegisterHandler function.
func (b *Builder) CreateMetricsFactory(namespace string) (metrics.Factory, error) {
	if b.Backend == "prometheus" {
		metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: namespace, Tags: nil})
		b.handler = promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true})
		return metricsFactory, nil
	}
	if b.Backend == "expvar" {
		metricsFactory := jexpvar.NewFactory(10).Namespace(metrics.NSOptions{Name: namespace, Tags: nil})
		b.handler = expvar.Handler()
		return metricsFactory, nil
	}
	if b.Backend == "none" || b.Backend == "" {
		return metrics.NullFactory, nil
	}
	return nil, errUnknownBackend
}

// Handler returns an http.Handler for the metrics endpoint.
func (b *Builder) Handler() http.Handler {
	return b.handler
}
