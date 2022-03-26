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
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	"github.com/Clymene-project/Clymene/pkg/multierror"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
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

type LogWriter struct {
	logger        *zap.Logger
	reporter      logproto.GatewayClient
	writerMetrics WriterMetrics
	marshaller    kafka.Marshaller
}

type LogWriterParams struct {
	Conn          *grpc.ClientConn
	Logger        *zap.Logger
	MetricFactory metrics.Factory
	Marshaller    kafka.Marshaller
}

func NewLogWriter(p *LogWriterParams) (*LogWriter, error) {
	writeMetrics := WriterMetrics{
		WrittenSuccess: p.MetricFactory.Counter(metrics.Options{Name: "gateway_grpc_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: p.MetricFactory.Counter(metrics.Options{Name: "gateway_grpc_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	return &LogWriter{
		reporter:      logproto.NewGatewayClient(p.Conn),
		logger:        p.Logger,
		writerMetrics: writeMetrics,
		marshaller:    p.Marshaller,
	}, nil
}
func (m *LogWriter) Writelog(ctx context.Context, tenantID string, batch logstore.Batch) (int, int64, int64, error) {
	var bufBytes int64
	var entriesCount64 int64
	var errs []error

	producerMessage := &client.ProducerBatch{TenantID: tenantID, Batch: *batch.(*client.Batch)}
	logsBytes, err := m.marshaller.MarshalLog(producerMessage)
	if err != nil {
		m.writerMetrics.WrittenFailure.Inc(1)
		errs = append(errs, err)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	}
	_, err = m.reporter.TransferBatch(ctx, &logproto.Batch{Batch: logsBytes})
	if err != nil {
		m.writerMetrics.WrittenFailure.Inc(1)
		errs = append(errs, err)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	} else {
		m.writerMetrics.WrittenSuccess.Inc(1)
	}
	return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
}
