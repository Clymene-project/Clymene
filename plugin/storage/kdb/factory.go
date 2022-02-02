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

package kdb

import (
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// https://code.kx.com/q/

type Factory struct {
}

func (f Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	//TODO implement me
	panic("implement me")
}

func (f Factory) CreateWriter() (metricstore.Writer, error) {
	//TODO implement me
	panic("implement me")
}

func NewFactory() *Factory {
	return &Factory{}
}
