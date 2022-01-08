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

package dbmodel

import (
	"github.com/Clymene-project/Clymene/prompb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"
	"time"
)

type Converter struct {
	Logger *zap.Logger
}

func (c *Converter) ConvertTsToPoint(metric prompb.TimeSeries) *write.Point {
	tags := make(map[string]string)
	for _, label := range metric.Labels {
		tags[label.Name] = label.Value
	}
	field := make(map[string]interface{})
	var timestamp int64
	for _, sample := range metric.Samples {
		timestamp = sample.Timestamp
		field["value"] = sample.Value
	}
	return influxdb2.NewPoint(tags["__name__"], tags, field, c.timestampMsToTime(timestamp))
}

func (c *Converter) timestampMsToTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}
