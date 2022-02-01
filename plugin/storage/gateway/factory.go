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
	"github.com/Clymene-project/Clymene/plugin/storage/gateway/http"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Factory struct {
	client  Client
	Options Options

	metricsFactory metrics.Factory
	logger         *zap.Logger
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "gateway"))
	//var client Client
	if f.Options.ServiceType == "grpc" {
		client, err := grpc.NewClient(f.Options.grpcOptions, f.metricsFactory, logger)
		if err != nil {
			return fmt.Errorf("failed to create gRPC connect: %w", err)
		}
		f.client = client
	} else {
		client, err := http.NewClient(f.Options.httpOptions, f.metricsFactory, logger)
		if err != nil {
			return fmt.Errorf("failed to create http connect: %w", err)
		}
		f.client = client
	}
	return nil
}

func (f *Factory) CreateWriter() (metricstore.Writer, error) {
	return createMetricWriter(f.client)
}

func NewFactory() *Factory {
	return &Factory{}
}

func createMetricWriter(client Client) (metricstore.Writer, error) {
	return client.CreateWriter()
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.Options.AddFlags(flagSet)
}

func (f *Factory) InitFromViper(v *viper.Viper) {
	f.Options.InitFromViper(v)
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromOptions(o Options) {
	f.Options = o
}
