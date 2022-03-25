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
	"github.com/Clymene-project/Clymene/cmd/promtail-ingester/app"
	"github.com/Clymene-project/Clymene/cmd/promtail-ingester/app/builder"
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
	ClymeneIngesterName = "promtail-ingester"
)

func main() {
	svc := flags.NewService(ports.IngesterAdminHTTP)
	svc.NoStorage = true
	version.Set(Version, BuildTime)

	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}

	v := viper.New()
	command := &cobra.Command{
		Use:   ClymeneIngesterName,
		Short: ClymeneIngesterName + " consumes from Kafka and send to db.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			logger.Info("start....", zap.String("component name", ClymeneIngesterName))
			logger.Info("build info", zap.String("version", Version), zap.String("build_time", BuildTime))

			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "clymene"})
			metricsFactory := baseFactory.Namespace(metrics.NSOptions{Name: "ingester"})

			storageFactory.InitFromViper(v)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}

			logWriter, err := storageFactory.CreateLogWriter()
			if err != nil {
				logger.Fatal("Failed to create metric writer", zap.Error(err))
			}

			options := app.Options{}
			options.InitFromViper(v) // default encode is protobuf
			consumer, err := builder.CreateConsumer(
				logger.With(zap.String("component", "consumer")),
				metricsFactory,
				logWriter,
				options,
			)
			if err != nil {
				logger.Fatal("Unable to create consumer", zap.Error(err))
			}
			consumer.Start()

			svc.RunAndThen(func() {
				if err := options.TLS.Close(); err != nil {
					logger.Error("Failed to close TLS certificates watcher", zap.Error(err))
				}
				if err = consumer.Close(); err != nil {
					logger.Error("Failed to close consumer", zap.Error(err))
				}
				if closer, ok := logWriter.(io.Closer); ok {
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
