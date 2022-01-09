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
	"go.uber.org/zap"
)

type MetricWriter struct {
	logger    *zap.Logger
	client    api.WriteAPI
	converter dbmodel.Converter
}

// WriteMetric is Asynchronous writer, and it is bulk insert
func (m *MetricWriter) WriteMetric(metrics []prompb.TimeSeries) error {
	for _, metric := range metrics {
		m.client.WritePoint(m.converter.ConvertTsToPoint(metric))
	}
	return nil
}

func NewMetricWriter(l *zap.Logger, client influxdb2.Client, org string, bucket string) *MetricWriter {
	return &MetricWriter{
		logger: l,
		client: client.WriteAPI(org, bucket),
		converter: dbmodel.Converter{
			Logger: l,
		},
	}
}
