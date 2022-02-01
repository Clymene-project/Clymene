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

package grpc

import (
	"context"
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MetricWriter struct {
	logger   *zap.Logger
	reporter prompb.ClymeneServiceClient
}

func (m *MetricWriter) WriteMetric(metric []prompb.TimeSeries) error {
	_, err := m.reporter.RequestMetrics(context.Background(), &prompb.WriteRequest{
		Timeseries: metric,
	})
	if err != nil {
		return err
	}
	return nil
}

type MetricWriterParams struct {
	Conn   *grpc.ClientConn
	Logger *zap.Logger
}

func NewMetricWriter(p *MetricWriterParams) (*MetricWriter, error) {
	return &MetricWriter{
		reporter: prompb.NewClymeneServiceClient(p.Conn),
		logger:   p.Logger,
	}, nil
}
