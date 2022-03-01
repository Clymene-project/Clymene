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
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

type MetricWriter struct {
	logger        *zap.Logger
	reporter      prompb.ClymeneServiceClient
	writerMetrics WriterMetrics
}

func (m *MetricWriter) WriteMetric(metric []prompb.TimeSeries) error {
	_, err := m.reporter.RequestMetrics(context.Background(), &prompb.WriteRequest{
		Timeseries: metric,
	})
	if err != nil {
		m.writerMetrics.WrittenFailure.Inc(1)
		return err
	}
	m.writerMetrics.WrittenSuccess.Inc(1)
	return nil
}

type MetricWriterParams struct {
	Conn          *grpc.ClientConn
	Logger        *zap.Logger
	MetricFactory metrics.Factory
}

func NewMetricWriter(p *MetricWriterParams) (*MetricWriter, error) {
	writeMetrics := WriterMetrics{
		WrittenSuccess: p.MetricFactory.Counter(metrics.Options{Name: "gateway_grpc_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: p.MetricFactory.Counter(metrics.Options{Name: "gateway_grpc_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	return &MetricWriter{
		reporter:      prompb.NewClymeneServiceClient(p.Conn),
		logger:        p.Logger,
		writerMetrics: writeMetrics,
	}, nil
}
