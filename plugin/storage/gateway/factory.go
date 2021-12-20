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

package gateway

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/gateway/grpc"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	grpc_conn "google.golang.org/grpc"
)

type Factory struct {
	Options *Options

	metricsFactory metrics.Factory
	logger         *zap.Logger
	conn           *grpc_conn.ClientConn
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger

	conn, err := f.Options.ConnBuilder.CreateConnection(logger)
	if err != nil {
		return fmt.Errorf("failed to create gRPC connect: %w", err)
	}
	f.conn = conn
	return nil
}

func (f *Factory) CreateWriter() (metricstore.Writer, error) {
	return createMetricWriter(f.logger, f.conn)
}

func NewFactory() *Factory {
	return &Factory{Options: NewOptions()}
}

func createMetricWriter(logger *zap.Logger, conn *grpc_conn.ClientConn) (metricstore.Writer, error) {
	return NewMetricWriter(&MetricWriterParams{Logger: logger, conn: conn})
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	grpc.AddFlags(flagSet)
}

func (f *Factory) InitFromViper(v *viper.Viper) {
	f.Options.InitFromViper(v)
}