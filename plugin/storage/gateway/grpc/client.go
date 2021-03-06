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
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	conn          *grpc.ClientConn
	logger        *zap.Logger
	metricFactory metrics.Factory
}

func (c *Client) CreateLogWriter() (logstore.Writer, error) {
	return NewLogWriter(&LogWriterParams{
		Conn:          c.conn,
		Logger:        c.logger,
		MetricFactory: c.metricFactory,
		Marshaller:    kafka.NewJSONMarshaller(),
	})
}

func (c *Client) CreateMetricWriter() (metricstore.Writer, error) {
	return NewMetricWriter(&MetricWriterParams{
		Conn:          c.conn,
		Logger:        c.logger,
		MetricFactory: c.metricFactory,
	})
}

func NewClient(options Options, metricFactory metrics.Factory, logger *zap.Logger) (*Client, error) {
	conn, err := options.CreateConnection(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connect: %w", err)
	}
	return &Client{
		conn:          conn,
		logger:        logger,
		metricFactory: metricFactory,
	}, nil
}
