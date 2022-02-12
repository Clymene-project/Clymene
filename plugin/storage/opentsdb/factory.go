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

package opentsdb

import (
	"flag"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/dryrun"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/http"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/socket"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// The opentsdb factory was developed based on opentsdb's tcollector.

type Factory struct {
	options Options

	metricsFactory metrics.Factory
	logger         *zap.Logger

	client Client
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	//TODO implement me
	panic("not supported")
}

func NewFactory() *Factory {
	return &Factory{}
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

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "opentsdb"))

	// hosts init
	var hosts []http.Hosts
	if f.options.hosts != "" {
		for _, host := range strings.Split(f.options.hosts, ",") {
			hostsplit := strings.Split(host, ":")
			var h http.Hosts
			var err error
			h.Host = hostsplit[0]
			h.Port, err = strconv.Atoi(hostsplit[1])
			// check host:port format
			if err != nil {
				f.logger.Panic("failed to host init", zap.Error(err))
			}
			hosts = append(hosts, h)
		}
	} else {
		hosts = append(hosts, http.Hosts{Host: f.options.host, Port: f.options.port})
	}

	// prompb.TimeSeries to json
	converter := &dbmodel.Converter{
		MaxTags: f.options.maxTags,
	}

	if f.options.http {
		// http, https
		f.client = http.NewClient(&http.Options{
			Hosts:        hosts,
			HttpPassword: f.options.httpPassword,
			HttpUsername: f.options.httpUsername,
			SSL:          f.options.ssl,
			HttpApiPath:  f.options.httpApiPath,
			Factory:      f.metricsFactory,
			Timeout:      f.options.timeout,
			MaxChunk:     f.options.maxChunk,
		}, converter, f.logger)
	} else if f.options.dryRun {
		// dryrun
		f.client = dryrun.NewClient(
			f.options.maxChunk,
			converter,
			f.logger,
		)
	} else {
		//	socket
		f.client = socket.NewClient(&socket.Options{
			Hosts: hosts,
		}, converter, f.logger)
	}

	return nil
}

func (f Factory) CreateMetricWriter() (metricstore.Writer, error) {
	return NewMetricWriter(f.logger, f.client), nil
}
