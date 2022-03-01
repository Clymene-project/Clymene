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

package opentsdb

import (
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
)

// The opentsdb factory was developed based on opentsdb's tcollector.

type Writer struct {
	client Client
	l      *zap.Logger
}

func NewMetricWriter(
	l *zap.Logger,
	c Client,
) *Writer {
	return &Writer{
		client: c,
		l:      l,
	}
}

func (w *Writer) WriteMetric(metrics []prompb.TimeSeries) error {
	return w.client.SendData(metrics)
}
