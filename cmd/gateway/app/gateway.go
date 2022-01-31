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
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/gateway/app/handler"
	"github.com/Clymene-project/Clymene/cmd/gateway/app/server"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"time"
)

type Gateway struct {
	logger         *zap.Logger
	metricsFactory metrics.Factory
	metricWriter   metricstore.Writer

	metricsHandler *handler.GRPCHandler

	grpcServer               *grpc.Server
	tlsGRPCCertWatcherCloser io.Closer
	httpServer               *http.Server
	tlsHTTPCertWatcherCloser io.Closer
}

type GatewayParams struct {
	Logger        *zap.Logger
	MetricFactory metrics.Factory
	MetricWriter  metricstore.Writer
}

func New(params *GatewayParams) *Gateway {
	return &Gateway{
		logger:         params.Logger,
		metricsFactory: params.MetricFactory,
		metricWriter:   params.MetricWriter,
	}
}

func (g *Gateway) Start(opt *GatewayOptions) error {
	grpcServer, err := server.StartGRPCServer(&server.GRPCServerParams{
		HostPort:      opt.gatewayGRPCHostPort,
		MetricHandler: handler.NewGRPCHandler(g.logger, g.metricWriter),
		TLSConfig:     opt.TLSGRPC,
		Logger:        g.logger,
	})
	if err != nil {
		return fmt.Errorf("could not start gRPC gateway %w", err)
	}
	httpServer, err := server.StartHTTPServer(&server.HTTPServerParams{
		HostPort:      opt.gatewayHTTPHostPort,
		MetricHandler: handler.NewHTTPHandler(g.logger, g.metricWriter),
		TLSConfig:     opt.TLSHTTP,
		Logger:        g.logger,
	})
	if err != nil {
		return fmt.Errorf("could not start HTTP gateway %w", err)
	}
	g.grpcServer = grpcServer
	g.httpServer = httpServer

	g.tlsGRPCCertWatcherCloser = &opt.TLSGRPC
	g.tlsHTTPCertWatcherCloser = &opt.TLSHTTP

	return nil
}

func (g *Gateway) Close() error {
	// gRPC server
	if g.grpcServer != nil {
		g.grpcServer.GracefulStop()
	}
	// HTTP server
	if g.httpServer != nil {

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := g.httpServer.Shutdown(timeout); err != nil {
			g.logger.Fatal("failed to stop the main HTTP server", zap.Error(err))
		}
		defer cancel()
	}

	// watchers actually never return errors from Close
	_ = g.tlsGRPCCertWatcherCloser.Close()
	_ = g.tlsHTTPCertWatcherCloser.Close()
	return nil
}
