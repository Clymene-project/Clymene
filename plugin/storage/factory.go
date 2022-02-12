// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/multierror"
	"github.com/Clymene-project/Clymene/plugin"
	"github.com/Clymene-project/Clymene/plugin/storage/cortex"
	"github.com/Clymene-project/Clymene/plugin/storage/es"
	"github.com/Clymene-project/Clymene/plugin/storage/gateway"
	"github.com/Clymene-project/Clymene/plugin/storage/influxdb"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/plugin/storage/kdb"
	"github.com/Clymene-project/Clymene/plugin/storage/loki"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb"
	"github.com/Clymene-project/Clymene/plugin/storage/prometheus"
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine"
	"github.com/Clymene-project/Clymene/storage"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"io"

	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

const (
	influxDbStorageType      = "influxdb"
	elasticsearchStorageType = "elasticsearch"
	prometheusStorageType    = "prometheus"
	kafkaStorageType         = "kafka"
	gatewayStorageType       = "gateway"
	cortexStorageType        = "cortex"
	kdbStorageType           = "kdb"
	opentsdbStorageType      = "opentsdb"
	tdengineStorageType      = "tdengine"
	lokiStorageType          = "loki"

	tsStorageType = "ts-storage-type"
)

// AllStorageTypes defines all available storage backends
var AllStorageTypes = []string{influxDbStorageType, elasticsearchStorageType, prometheusStorageType, kafkaStorageType,
	gatewayStorageType, cortexStorageType, kdbStorageType, opentsdbStorageType, tdengineStorageType, lokiStorageType}

// Factory implements storage.Factory interface as a meta-factory for storage components.
type Factory struct {
	FactoryConfig
	metricsFactory metrics.Factory
	factories      map[string]storage.Factory
}

// NewFactory creates the meta-factory.
func NewFactory(config FactoryConfig) (*Factory, error) {
	f := &Factory{FactoryConfig: config}
	uniqueTypes := map[string]struct{}{
		f.ReaderType:              {},
		f.DependenciesStorageType: {},
	}
	for _, storageType := range f.WriterTypes {
		uniqueTypes[storageType] = struct{}{}
	}
	f.factories = make(map[string]storage.Factory)
	for t := range uniqueTypes {
		ff, err := f.getFactoryOfType(t)
		if err != nil {
			return nil, err
		}
		f.factories[t] = ff
	}
	return f, nil
}

func (f *Factory) getFactoryOfType(factoryType string) (storage.Factory, error) {
	switch factoryType {
	case elasticsearchStorageType:
		return es.NewFactory(), nil
	case prometheusStorageType:
		return prometheus.NewFactory(), nil
	case cortexStorageType:
		return cortex.NewFactory(), nil
	case kafkaStorageType:
		return kafka.NewFactory(), nil
	case influxDbStorageType:
		return influxdb.NewFactory(), nil
	case gatewayStorageType:
		return gateway.NewFactory(), nil
	case opentsdbStorageType:
		return opentsdb.NewFactory(), nil
	case kdbStorageType:
		return kdb.NewFactory(), nil
	case tdengineStorageType:
		return tdengine.NewFactory(), nil
	case lokiStorageType:
		return loki.NewFactory(), nil
	default:
		return nil, fmt.Errorf("unknown storage type %s. Valid types are %v", factoryType, AllStorageTypes)
	}
}

// Initialize implements storage.Factory.
func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory = metricsFactory
	for _, factory := range f.factories {
		if err := factory.Initialize(metricsFactory, logger); err != nil {
			return err
		}
	}
	f.publishOpts()

	return nil
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.AddFlags(flagSet)
		}
	}
}

func (f *Factory) AddPipelineFlags(flagSet *flag.FlagSet) {
	f.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.InitFromViper(v)
		}
	}
}

var _ io.Closer = (*Factory)(nil)

// Close closes the resources held by the factory
func (f *Factory) Close() error {
	var errs []error
	for _, storageType := range f.WriterTypes {
		if factory, ok := f.factories[storageType]; ok {
			if closer, ok := factory.(io.Closer); ok {
				err := closer.Close()
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return multierror.Wrap(errs)
}

func (f *Factory) publishOpts() {
	internalFactory := f.metricsFactory.Namespace(metrics.NSOptions{Name: "internal"})
	internalFactory.Gauge(metrics.Options{Name: tsStorageType + "-" + f.ReaderType}).
		Update(1)
}

func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
	var writers []metricstore.Writer
	for _, storageType := range f.WriterTypes {
		factory, ok := f.factories[storageType]
		if !ok {
			return nil, fmt.Errorf("no %s backend registered for metric store", storageType)
		}
		writer, err := factory.CreateMetricWriter()
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	var Writer metricstore.Writer
	if len(f.WriterTypes) == 1 {
		Writer = writers[0]
	} else {
		Writer = metricstore.NewCompositeWriter(writers...)
	}
	return Writer, nil
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	var writers []logstore.Writer
	for _, storageType := range f.WriterTypes {
		factory, ok := f.factories[storageType]
		if !ok {
			return nil, fmt.Errorf("no %s backend registered for metric store", storageType)
		}
		writer, err := factory.CreateLogWriter()
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	var Writer logstore.Writer
	if len(f.WriterTypes) == 1 {
		Writer = writers[0]
	} else {
		Writer = logstore.NewCompositeWriter(writers...)
	}
	return Writer, nil
}
