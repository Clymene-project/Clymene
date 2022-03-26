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

package http

import (
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Client struct {
	options        Options
	logger         *zap.Logger
	metricsFactory metrics.Factory
}

func (c *Client) CreateLogWriter() (logstore.Writer, error) {
	return NewLogWriter(c.logger, c.metricsFactory, c.options, kafka.NewJSONMarshaller()), nil
}

func (c *Client) CreateMetricWriter() (metricstore.Writer, error) {
	return NewMetricWriter(c.logger, c.metricsFactory, c.options, kafka.NewProtobufMarshaller()), nil
}

func NewClient(options Options, factory metrics.Factory, logger *zap.Logger) (*Client, error) {
	return &Client{
		options:        options,
		metricsFactory: factory,
		logger:         logger,
	}, nil
}
