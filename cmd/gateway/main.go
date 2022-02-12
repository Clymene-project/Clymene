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

package main

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/docs"
	"github.com/Clymene-project/Clymene/cmd/flags"
	"github.com/Clymene-project/Clymene/cmd/gateway/app"
	"github.com/Clymene-project/Clymene/pkg/config"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/plugin/storage"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

var (
	// Version is set during binary building (git revision number)
	Version string
	// BuildTime is set during binary building
	BuildTime string
)

const (
	ClymeneGatewayName = "Clymene-gateway"
)

func main() {
	svc := flags.NewService(ports.GatewayAdminHTTP)
	svc.NoStorage = true
	version.Set(Version, BuildTime)

	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}

	v := viper.New()
	command := &cobra.Command{
		Use:   ClymeneGatewayName,
		Short: ClymeneGatewayName + " can receive data through gRPC.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			logger.Info("start....", zap.String("component name", ClymeneGatewayName))
			logger.Info("build info", zap.String("version", Version), zap.String("build_time", BuildTime))

			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "clymene"})
			metricsFactory := baseFactory.Namespace(metrics.NSOptions{Name: "gateway"})

			storageFactory.InitFromViper(v)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}

			metricWriter, err := storageFactory.CreateMetricWriter()
			if err != nil {
				logger.Fatal("Failed to create metric writer", zap.Error(err))
			}
			gatewayOpt := new(app.GatewayOptions).InitFromViper(v)

			gateway := app.New(&app.GatewayParams{
				Logger:        logger,
				MetricFactory: metricsFactory,
				MetricWriter:  metricWriter,
			})

			if err := gateway.Start(gatewayOpt); err != nil {
				log.Fatal(err)
			}

			svc.RunAndThen(func() {
				if err = gateway.Close(); err != nil {
					logger.Error("Failed to close gateway", zap.Error(err))
				}
				if closer, ok := metricWriter.(io.Closer); ok {
					err := closer.Close()
					if err != nil {
						logger.Error("Failed to close metrics writer", zap.Error(err))
					}
				}
				if err := storageFactory.Close(); err != nil {
					logger.Error("Failed to close storage factory", zap.Error(err))
				}
			})
			return nil
		},
	}

	command.AddCommand(version.Command())
	command.AddCommand(docs.Command(v))

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		storageFactory.AddPipelineFlags,
		app.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
