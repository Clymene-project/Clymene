// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"time"

	"github.com/uber/jaeger-lib/metrics"
)

// WriteLogs is a collection of metrics for write operations.
type WriteLogs struct {
	Attempts   metrics.Counter `metric:"attempts"`
	Inserts    metrics.Counter `metric:"inserts"`
	Errors     metrics.Counter `metric:"errors"`
	LatencyOk  metrics.Timer   `metric:"latency-ok"`
	LatencyErr metrics.Timer   `metric:"latency-err"`
}

// NewWritelogs takes a metrics scope and creates a metrics struct
func NewWritelogs(factory metrics.Factory, tableName string) *WriteLogs {
	t := &WriteLogs{}
	metrics.Init(t, factory.Namespace(metrics.NSOptions{Name: tableName, Tags: nil}), nil)
	return t
}

// Emit will record success or failure counts and latency metrics depending on the passed error.
func (t *WriteLogs) Emit(err error, latency time.Duration) {
	t.Attempts.Inc(1)
	if err != nil {
		t.LatencyErr.Record(latency)
		t.Errors.Inc(1)
	} else {
		t.LatencyOk.Record(latency)
		t.Inserts.Inc(1)
	}
}
