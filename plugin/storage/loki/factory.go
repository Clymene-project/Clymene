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

package loki

import (
	"flag"
	lokiflag "github.com/Clymene-project/Clymene/pkg/lokiutil/flagext"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Factory struct {
	Options Options

	client *http.Client

	metricsFactory metrics.Factory
	logger         *zap.Logger
	externalLabels lokiflag.LabelSet
	URL            *url.URL
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	return NewLogWriter(f.client, f.URL, f.externalLabels, f.logger), nil
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	if f.Options.URL.URL == nil {
		return errors.New("client needs target URL")
	}

	f.URL = f.Options.URL.URL
	err := f.Options.Client.Validate()
	if err != nil {
		return err
	}
	f.client, err = config.NewClientFromConfig(f.Options.Client, "promtail", config.WithHTTP2Disabled())
	if err != nil {
		return err
	}
	f.client.Timeout = f.Options.Timeout
	return nil
}

func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
	//TODO implement me
	panic("not supported")
}

func NewFactory() *Factory {
	return &Factory{}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.Options.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.Options.InitFromViper(v)
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromOptions(o Options) {
	f.Options = o
}
