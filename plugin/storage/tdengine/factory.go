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

package tdengine

import (
	"database/sql"
	"flag"
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine/db"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Factory implements storage.Factory
type Factory struct {
	options        Options
	logger         *zap.Logger
	metricsFactory metrics.Factory
	tdEngine       *sql.DB
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	//TODO implement me
	panic("not supported")
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "tdengine"))
	logger.Info("TDengine factory", zap.String("url", f.options.hostName))

	// TDengine session
	tdEngine, err := db.PrepareConnection(&db.TaosConnConfig{
		User:       f.options.user,
		Password:   f.options.password,
		HostName:   f.options.hostName,
		ServerPort: f.options.serverPort,
		DbName:     f.options.dbName,
	})
	if err != nil {
		logger.Panic("connection Error", zap.Error(err))
	}
	f.tdEngine = tdEngine
	return nil
}

func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
	return NewMetricWriter(f.tdEngine, f.options.maxSQLLength, f.metricsFactory, f.logger), nil
}

// NewFactory creates a new Factory.
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
