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

package storage

import (
	"github.com/Clymene-project/Clymene/storage/metricstore"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Factory defines an interface for a factory that can create implementations of different storage components.
// Implementations are also encouraged to implement plugin.Configurable interface.
//
// See also
//
// plugin.Configurable
type Factory interface {
	// Initialize performs internal initialization of the factory, such as opening connections to the backend store.
	// It is called after all configuration of the factory itself has been done.
	Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error

	CreateWriter() (metricstore.Writer, error)
}
