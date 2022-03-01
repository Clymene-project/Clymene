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

package influxdb

import (
	"github.com/Clymene-project/Clymene/plugin/storage/influxdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

type MetricWriter struct {
	logger       *zap.Logger
	client       api.WriteAPI
	converter    dbmodel.Converter
	writeMetrics WriterMetrics
}

// WriteMetric is Asynchronous writer, and it is bulk insert
func (m *MetricWriter) WriteMetric(metrics []prompb.TimeSeries) error {
	for _, metric := range metrics {
		m.client.WritePoint(m.converter.ConvertTsToPoint(metric))
	}
	m.writeMetrics.WrittenSuccess.Inc(1)
	return nil
}

func NewMetricWriter(l *zap.Logger, client influxdb2.Client, org string, bucket string, metricFactory metrics.Factory) *MetricWriter {
	writeMetrics := WriterMetrics{
		WrittenSuccess: metricFactory.Counter(metrics.Options{Name: "influxdb_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: metricFactory.Counter(metrics.Options{Name: "influxdb_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	writerClient := client.WriteAPI(org, bucket)
	go func() {
		for range writerClient.Errors() {
			writeMetrics.WrittenFailure.Inc(1)
		}
	}()
	return &MetricWriter{
		logger: l,
		client: writerClient,
		converter: dbmodel.Converter{
			Logger: l,
		},
		writeMetrics: writeMetrics,
	}
}
