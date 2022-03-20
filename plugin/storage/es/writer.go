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
	"github.com/Clymene-project/Clymene/pkg/multierror"
	"github.com/Clymene-project/Clymene/plugin/storage/es/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"go.uber.org/zap"
)

const (
	metricType = "metric"
	logType    = "log"
)

type Writer struct {
	logger    *zap.Logger
	client    es.Client
	index     string
	converter dbmodel.Converter
}

func (s *Writer) Writelog(ctx context.Context, tenantID string, batch logstore.Batch) (int, int64, int64, error) {
	var statusCode int
	var bufBytes int64
	var entriesCount64 int64
	var errs []error
	streams, entriesCount := batch.CreatePushRequest()
	entriesCount64 = int64(entriesCount)
	bufBytes = int64(len(streams.Streams)) // Put the length value of Streams for data verification
	for _, stream := range streams.Streams {
		labels, err := s.converter.ConvertStringToLabel(stream.Labels)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, entry := range stream.Entries {
			logs, err := s.converter.ConvertLogsToJSON(tenantID, labels, entry, stream.Hash)
			s.writelog(logs)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return statusCode, bufBytes, entriesCount64, multierror.Wrap(errs)
}

// bulk insert
func (s *Writer) writelog(logs *map[string]interface{}) {
	s.client.Index().Index(s.index).Type(logType).BodyJson(&logs).Add()
}

// WriterParams holds constructor parameters for NewMetricWriter
type WriterParams struct {
	Logger      *zap.Logger
	Client      es.Client
	IndexPrefix string
	Archive     bool
	Index       string
}

// NewMetricWriter creates a new MetricWriter for use
func NewMetricWriter(p WriterParams) *Writer {
	prefix := ""
	if p.IndexPrefix != "" {
		prefix = p.IndexPrefix + "-"
	}
	return &Writer{
		client:    p.Client,
		logger:    p.Logger,
		index:     prefix + p.Index,
		converter: dbmodel.Converter{},
	}
}

func (s *Writer) WriteMetric(metrics []prompb.TimeSeries) error {
	for _, metric := range metrics {
		jsonTimeSeries := s.converter.ConvertTsToJSON(metric)
		s.writeMetric(&jsonTimeSeries)
	}
	return nil
}

// bulk insert
func (s *Writer) writeMetric(metric *map[string]interface{}) {
	s.client.Index().Index(s.index).Type(metricType).BodyJson(&metric).Add()
}

// NewLogWriter creates a new MetricWriter for use
func NewLogWriter(p WriterParams) *Writer {
	prefix := ""
	if p.IndexPrefix != "" {
		prefix = p.IndexPrefix + "-"
	}
	return &Writer{
		client:    p.Client,
		logger:    p.Logger,
		index:     prefix + p.Index,
		converter: dbmodel.Converter{},
	}
}
