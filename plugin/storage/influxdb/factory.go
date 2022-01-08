/*
 * Copyright (c) 2021 The Clymene Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package influxdb

import (
	"flag"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Factory struct {
	options Options

	metricsFactory metrics.Factory
	logger         *zap.Logger
	client         influxdb2.Client
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "influxdb"))
	if f.options.TLS.Enabled {
		tls, err := f.options.TLS.Config(f.logger)
		if err != nil {
			return err
		}
		f.options.Options.SetTLSConfig(tls)
	}
	f.options.SetHTTPClient(
		&http.Client{
			Transport: newLatencyTransport(http.DefaultTransport, metricsFactory),
			Timeout:   time.Second * time.Duration(f.options.HTTPRequestTimeout()),
		})
	err := f.options.checkNecessaryOptions()
	if err != nil {
		return err
	}
	f.client = influxdb2.NewClientWithOptions(f.options.url, f.options.token, &f.options.Options)
	return nil
}

func (f *Factory) CreateWriter() (metricstore.Writer, error) {
	return NewMetricWriter(f.logger, f.client, f.options.org, f.options.bucket), nil
}

func NewFactory() *Factory {
	return &Factory{
		options: Options{},
	}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.options.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.options.InitFromViper(v)
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromOptions(o Options) {
	f.options = o
}

type latencyTransport struct {
	transport http.RoundTripper
	latency   metrics.Timer
	errors    metrics.Counter
}

func (l *latencyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	now := time.Now()
	resp, err := l.transport.RoundTrip(request)
	if err != nil {
		l.errors.Inc(1)
		return resp, err
	}
	l.latency.Record(time.Since(now))
	return resp, err
}

func newLatencyTransport(t http.RoundTripper, f metrics.Factory) http.RoundTripper {
	m := f.Namespace(metrics.NSOptions{Name: "influxdb", Tags: nil})
	return &latencyTransport{
		transport: t,
		latency:   m.Timer(metrics.TimerOptions{Name: "latency", Tags: nil}),
		errors:    m.Counter(metrics.Options{Name: "errors", Tags: nil}),
	}
}
