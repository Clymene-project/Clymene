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
	"github.com/bourbonkk/Clymene/cmd/agent/app/reporter/grpc"
	"github.com/bourbonkk/Clymene/cmd/docs"
	"github.com/bourbonkk/Clymene/cmd/flags"
	"github.com/bourbonkk/Clymene/pkg/config"
	"github.com/bourbonkk/Clymene/pkg/version"
	"github.com/bourbonkk/Clymene/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	svc := flags.NewService(ports.AgentAdminHTTP)
	svc.NoStorage = true

	v := viper.New()
	var command = &cobra.Command{
		Use:   "clymene-agent",
		Short: "clymene agent is a local daemon program which scrapes metric data.",
		Long:  `clymene agent is a daemon program that runs on every cluster and scrapes metric data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger // shortcut
			logger.Info("starting")

			//grpcBuilder := grpc.NewConnBuilder().InitFromViper(v)
			//builders := map[reporter.Type]app.CollectorProxyBuilder{
			//	reporter.GRPC: app.GRPCCollectorProxyBuilder(grpcBuilder),
			//}
			//cp, err := app.CreateCollectorProxy(app.ProxyBuilderOptions{
			//	Options: *rOpts,
			//	Logger:  logger,
			//	Metrics: mFactory,
			//}, builders)
			//if err != nil {
			//	logger.Fatal("Could not create collector proxy", zap.Error(err))
			//}
			//
			//// TODO illustrate discovery service wiring
			//
			//builder := new(app.Builder).InitFromViper(v)
			//agent, err := builder.CreateAgent(cp, logger, mFactory)
			//if err != nil {
			//	return fmt.Errorf("unable to initialize Jaeger Agent: %w", err)
			//}
			//
			//logger.Info("Starting agent")
			//if err := agent.Run(); err != nil {
			//	return fmt.Errorf("failed to run the agent: %w", err)
			//}

			svc.RunAndThen(func() {
				//agent.Stop()
				//cp.Close()
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
		grpc.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
