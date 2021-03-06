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

package es

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/es"
	"github.com/Clymene-project/Clymene/pkg/es/config"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

const (
	primaryNamespace = "es"
	archiveNamespace = "es-archive"
)

// Factory implements storage.Factory for Elasticsearch backend.
type Factory struct {
	Options *Options

	metricsFactory metrics.Factory
	logger         *zap.Logger

	primaryConfig config.ClientBuilder
	primaryClient es.Client
	archiveConfig config.ClientBuilder
	archiveClient es.Client
}

// NewFactory creates a new Factory.
func NewFactory() *Factory {
	return &Factory{
		Options: NewOptions(primaryNamespace, archiveNamespace),
	}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.Options.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.Options.InitFromViper(v)
	f.primaryConfig = f.Options.GetPrimary()
	f.archiveConfig = f.Options.Get(archiveNamespace)
}

// InitFromOptions configures factory from Options struct.
func (f *Factory) InitFromOptions(o Options) {
	f.Options = &o
	f.primaryConfig = f.Options.GetPrimary()
	if cfg := f.Options.Get(archiveNamespace); cfg != nil {
		f.archiveConfig = cfg
	}
}

// Initialize implements storage.Factory
func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "elasticsearch"))

	primaryClient, err := f.primaryConfig.NewClient(logger, metricsFactory)
	if err != nil {
		return fmt.Errorf("failed to create primary Elasticsearch client: %w", err)
	}
	f.primaryClient = primaryClient
	if f.archiveConfig.IsStorageEnabled() {
		f.archiveClient, err = f.archiveConfig.NewClient(logger, metricsFactory)
		if err != nil {
			return fmt.Errorf("failed to create archive Elasticsearch client: %w", err)
		}
	}
	return nil
}

func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
	return createMetricWriter(f.logger, f.primaryClient, f.primaryConfig, f.Options.metricsIndex, false)
}

func createMetricWriter(logger *zap.Logger, client es.Client, cfg config.ClientBuilder, metricsIndex string, archive bool) (metricstore.Writer, error) {
	writer := NewMetricWriter(WriterParams{
		Client:      client,
		Logger:      logger,
		IndexPrefix: cfg.GetIndexPrefix(),
		Index:       metricsIndex,
		Archive:     archive,
	})

	return writer, nil
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	return createLogWriter(f.logger, f.primaryClient, f.primaryConfig, f.Options.logsIndex, false)
}

func createLogWriter(logger *zap.Logger, client es.Client, cfg config.ClientBuilder, logsIndex string, archive bool) (logstore.Writer, error) {
	writer := NewLogWriter(WriterParams{
		Client:      client,
		Logger:      logger,
		IndexPrefix: cfg.GetIndexPrefix(),
		Index:       logsIndex,
		Archive:     archive,
	})
	return writer, nil
}
