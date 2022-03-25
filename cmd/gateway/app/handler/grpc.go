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

package handler

import (
	"context"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"go.uber.org/zap"
)

// GRPCHandler implements gRPC CollectorService.
type GRPCHandler struct {
	logger       *zap.Logger
	metricWriter metricstore.Writer
	logWriter    logstore.Writer
}

func (g *GRPCHandler) RequestMetrics(c context.Context, r *prompb.WriteRequest) (*prompb.MetricsResponse, error) {
	err := g.metricWriter.WriteMetric(r.GetTimeseries())
	if err != nil {
		g.logger.Warn("Failed to create metric", zap.Error(err))
		return nil, err
	}
	return &prompb.MetricsResponse{}, nil
}

func NewGRPCHandler(logger *zap.Logger, metricWriter metricstore.Writer, logWriter logstore.Writer) *GRPCHandler {
	return &GRPCHandler{
		logger:       logger,
		metricWriter: metricWriter,
		logWriter:    logWriter,
	}
}
