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
	"github.com/Clymene-project/Clymene/cmd/promtail/app"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/pkg/config"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/plugin/storage"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
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
	ClymenePromtailName = "Clymene-promtail"
)

func main() {
	svc := flags.NewService(ports.PromtailAdminHTTP)
	svc.NoStorage = true
	version.Set(Version, BuildTime)

	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}

	v := viper.New()
	command := &cobra.Command{
		Use:   ClymenePromtailName,
		Short: ClymenePromtailName + " is a log collection agent.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			logger.Info("start....", zap.String("component name", ClymenePromtailName))
			logger.Info("build info", zap.String("version", Version), zap.String("build_time", BuildTime))

			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "clymene"})
			metricsFactory := baseFactory.Namespace(metrics.NSOptions{Name: "promtail"})

			storageFactory.InitFromViper(v)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}

			logWriter, err := storageFactory.CreateLogWriter()
			if err != nil {
				logger.Fatal("Failed to create log writer", zap.Error(err))
			}

			options := client.Options{}
			options.InitFromViper(v)
			promtail, err := app.New(&app.PromtailConfig{
				Options:        options,
				Reg:            prometheus.DefaultRegisterer,
				MetricsFactory: metricsFactory,
				LogWriter:      logWriter,
				Logger:         logger.With(zap.String("component", "promtail")),
			})
			if err != nil {
				logger.Fatal("Unable to create promtail", zap.Error(err))
			}

			if err := promtail.Run(); err != nil {
				logger.Error("error starting promtail", zap.Error(err))
			}
			svc.RunAndThen(func() {
				promtail.Shutdown()
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
		client.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
