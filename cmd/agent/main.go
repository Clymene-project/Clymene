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
	"github.com/Clymene-project/Clymene/cmd/status"
	"github.com/Clymene-project/Clymene/pkg/config"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/plugin/storage"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	svc := flags.NewService(ports.AgentAdminHTTP)

	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}

	v := viper.New()
	var command = &cobra.Command{
		Use:   "clymene-agent",
		Short: "clymene agent is a local daemon program which scrapes metric data.",
		Long:  `clymene agent is a daemon program that runs on every cluster and scrapes metric data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "clymene"})
			//metricsFactory := baseFactory.Namespace(metrics.NSOptions{Name: "collector"})

			storageFactory.InitFromViper(v)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
			}

			svc.RunAndThen(func() {
				//agent.Stop()
				//cp.Close()
			})
			return nil
		},
	}

	command.AddCommand(version.Command())
	command.AddCommand(docs.Command(v))
	command.AddCommand(status.Command(v, ports.AgentAdminHTTP))

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		storageFactory.AddPipelineFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
