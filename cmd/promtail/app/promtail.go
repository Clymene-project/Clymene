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
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/config"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/logentry/stages"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/server"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Promtail is the root struct for Promtail.
type Promtail struct {
	client         client.Client
	targetManagers *targets.TargetManagers
	server         server.Server
	logger         *zap.Logger
	reg            prometheus.Registerer

	stopped bool
	mtx     sync.Mutex
}

type PromtailConfig struct {
	Options
	Reg            prometheus.Registerer
	Logger         *zap.Logger
	MetricsFactory metrics.Factory
	LogWriter      logstore.Writer
}

// New makes a new Promtail.
func New(pc *PromtailConfig) (*Promtail, error) {
	promtail := &Promtail{
		logger: pc.Logger,
		reg:    pc.Reg,
	}
	cfg, err := config.LoadFile(pc.configFile, pc.Logger)
	if err != nil {
		pc.Logger.Panic("Unable to parse config")
	}

	if cfg.LimitConfig.ReadlineRateEnabled {
		stages.SetReadLineRateLimiter(cfg.LimitConfig.ReadlineRate, cfg.LimitConfig.ReadlineBurst, cfg.LimitConfig.ReadlineRateDrop)
	}
	if pc.dryRun {
		promtail.client, err = client.NewLogger(pc.Reg, promtail.logger, cfg.ClientConfigs...)
		if err != nil {
			return nil, err
		}
		cfg.PositionsConfig.ReadOnly = true
	} else {
		promtail.client, err = client.NewMulti(pc.Reg, promtail.logger, cfg.ClientConfigs...)
		if err != nil {
			return nil, err
		}
	}

	tms, err := targets.NewTargetManagers(promtail, promtail.reg, promtail.logger, cfg.PositionsConfig, promtail.client, cfg.ScrapeConfig, &cfg.TargetConfig)
	if err != nil {
		return nil, err
	}
	promtail.targetManagers = tms
	//server, err := server.New(cfg.ServerConfig, promtail.logger, tms, cfg.String())
	if err != nil {
		return nil, err
	}
	promtail.server = nil
	return promtail, nil
}

// Run the promtail; will block until a signal is received.
func (p *Promtail) Run() error {
	p.mtx.Lock()
	// if we stopped promtail before the server even started we can return without starting.
	if p.stopped {
		p.mtx.Unlock()
		return nil
	}
	p.mtx.Unlock() // unlock before blocking
	return p.server.Run()
}

// Client returns the underlying client Promtail uses to write to Loki.
func (p *Promtail) Client() client.Client {
	return p.client
}

// Shutdown the promtail.
func (p *Promtail) Shutdown() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.stopped = true
	if p.server != nil {
		p.server.Shutdown()
	}
	if p.targetManagers != nil {
		p.targetManagers.Stop()
	}
	// todo work out the stop.
	p.client.Stop()
}
