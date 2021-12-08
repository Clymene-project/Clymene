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

package app

import (
	"context"
	agent_config "github.com/Clymene-project/Clymene/cmd/agent/app/config"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery"
	sd_config "github.com/Clymene-project/Clymene/cmd/agent/app/discovery/config"
	"github.com/Clymene-project/Clymene/cmd/agent/app/scrape"
	"github.com/Clymene-project/Clymene/cmd/agent/app/server"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

type Agent struct {
	l *zap.Logger

	g *run.Group

	writer metricstore.Writer

	hServer    *http.Server
	hostPort   int
	configFile string

	reloaders []func(cfg *agent_config.Config) error
}

type AgentConfig struct {
	ConfigFile    string
	MetricWriter  metricstore.Writer
	Logger        *zap.Logger
	HttpPort      int
	MetricFactory metrics.Factory
}

// New constructs a new agent component
func New(config *AgentConfig) *Agent {
	var scrapeManager *scrape.Manager

	type closeOnce struct {
		C     chan struct{}
		once  sync.Once
		Close func()
	}
	var g run.Group
	// Wait until the server is ready to handle reloading.
	reloadReady := &closeOnce{
		C: make(chan struct{}),
	}
	reloadReady.Close = func() {
		reloadReady.once.Do(func() {
			close(reloadReady.C)
		})
	}

	ctxScrape, cancelScrape := context.WithCancel(context.Background())
	discoveryManagerScrape := discovery.NewManager(ctxScrape, config.Logger.With(zap.String("component", "discovery manager scrape")), discovery.Name("scrape"))

	scrapeManager = scrape.NewManager(config.Logger.With(zap.String("component", "scrape manager")), config.MetricWriter)

	reloaders := []func(cfg *agent_config.Config) error{
		scrapeManager.ApplyConfig,
		func(cfg *agent_config.Config) error {
			c := make(map[string]sd_config.ServiceDiscoveryConfig)
			for _, v := range cfg.ScrapeConfigs {
				c[v.JobName] = v.ServiceDiscoveryConfig
			}
			return discoveryManagerScrape.ApplyConfig(c)
		},
	}
	{
		// Scrape discovery manager.
		g.Add(
			func() error {
				err := discoveryManagerScrape.Run()
				config.Logger.Info("Scrape discovery manager stopped")
				return err
			},
			func(err error) {
				config.Logger.Info("Stopping scrape discovery manager...")
				cancelScrape()
			},
		)
	}
	{
		// Scrape manager
		g.Add(
			func() error {
				<-reloadReady.C
				err := scrapeManager.Run(discoveryManagerScrape.SyncCh())
				config.Logger.Info("Scrape manager stopped")
				return err
			},
			func(err error) {
				config.Logger.Info("Stopping scrape manager...")
				scrapeManager.Stop()
			},
		)
	}

	{
		// Initial config
		cancel := make(chan struct{})
		g.Add(
			func() error {
				if err := agent_config.ReloadConfig(config.ConfigFile, config.Logger, reloaders...); err != nil {
					return errors.Wrapf(err, "error loading config from %q", config.ConfigFile)
				}
				reloadReady.Close()
				config.Logger.Info("Server is ready to receive web requests.")
				<-cancel
				return nil
			},
			func(err error) {
				close(cancel)
			},
		)
	}

	return &Agent{
		l:          config.Logger,
		hostPort:   config.HttpPort,
		configFile: config.ConfigFile,
		g:          &g,
		reloaders:  reloaders,
	}
}

func (a *Agent) Run() error {
	// HTTP server for config reload
	if httpServer, err := server.StartHTTPServer(&server.HttpServerParams{
		HostPort:   ports.PortToHostPort(a.hostPort),
		Logger:     a.l,
		Reloader:   a.reloaders,
		ConfigFile: a.configFile,
	}); err != nil {
		return err
	} else {
		a.hServer = httpServer
	}
	if err := a.g.Run(); err != nil {
		a.l.Error("agent group run error", zap.Error(err))
		os.Exit(1)
	}
	return nil
}

func (a *Agent) ReloadConfig() error {
	if err := agent_config.ReloadConfig(a.configFile, a.l, a.reloaders...); err != nil {
		return errors.Wrapf(err, "error loading config from %q", a.configFile)
	}
	return nil
}

func (a *Agent) Stop() error {
	a.Close()
	os.Exit(1)
	return nil
}

func (a *Agent) Close() {
	if err := a.hServer.Close(); err != nil {
		a.l.Error("agent httpserver close error", zap.Error(err))
	}
}
