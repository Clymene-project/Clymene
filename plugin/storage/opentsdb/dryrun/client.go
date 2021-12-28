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

package dryrun

import (
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
)

type Client struct {
	maxChunk  int
	converter *dbmodel.Converter
	logger    *zap.Logger
}
type Options struct {
	maxChunk int
}

func (c *Client) SendData(metrics []prompb.TimeSeries) error {
	q := len(metrics) / c.maxChunk
	r := len(metrics) % c.maxChunk
	if r != 0 {
		q += 1
	}
	for i := 1; i <= q; i++ {
		var timeSeriesDiv []prompb.TimeSeries
		if i == 1 {
			timeSeriesDiv = metrics[:i*c.maxChunk]
		} else if i != q {
			timeSeriesDiv = metrics[(i-1)*c.maxChunk : i*c.maxChunk]
		} else {
			timeSeriesDiv = metrics[(i-1)*c.maxChunk:]
		}
		jsonTS, _ := c.converter.ConvertTsToOpenTSDBJSON(timeSeriesDiv)

		c.logger.Info("dryRun", zap.String("sendData", string(jsonTS)))
	}
	return nil
}

func NewClient(maxChunk int, converter *dbmodel.Converter, l *zap.Logger) *Client {
	return &Client{
		maxChunk:  maxChunk,
		converter: converter,
		logger:    l,
	}
}
