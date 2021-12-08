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

package es

import (
	"context"
	"github.com/Clymene-project/Clymene/pkg/es"
	"github.com/Clymene-project/Clymene/prompb"
	storageMetrics "github.com/Clymene-project/Clymene/storage/metricstore/metrics"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"strings"
)

type MetricWriter struct {
	client        es.Client
	logger        *zap.Logger
	writerMetrics metricWriterMetrics
}

type metricWriterMetrics struct {
	indexCreate *storageMetrics.WriteMetrics
}

// MetricWriterParams holds constructor parameters for NewMetricWriter
type MetricWriterParams struct {
	Client              es.Client
	Logger              *zap.Logger
	MetricsFactory      metrics.Factory
	IndexPrefix         string
	IndexDateLayout     string
	AllTagsAsFields     bool
	TagKeysAsFields     []string
	TagDotReplacement   string
	Archive             bool
	UseReadWriteAliases bool
}

// NewMetricWriter creates a new MetricWriter for use
func NewMetricWriter(p MetricWriterParams) *MetricWriter {
	return &MetricWriter{
		client: p.Client,
		logger: p.Logger,
		writerMetrics: metricWriterMetrics{
			indexCreate: storageMetrics.NewWriteMetrics(p.MetricsFactory, "index_create"),
		},
	}
}

// CreateTemplates creates index templates.
func (s *MetricWriter) CreateTemplates(metricTemplate, indexPrefix string) error {
	if indexPrefix != "" && !strings.HasSuffix(indexPrefix, "-") {
		indexPrefix += "-"
	}
	_, err := s.client.CreateTemplate(indexPrefix + "clymene-metrics").Body(metricTemplate).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricWriter) WriteMetric(metric []prompb.TimeSeries) error {
	//TODO implement me
	panic("implement me")
}
