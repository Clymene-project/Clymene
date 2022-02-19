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

package logstore

import (
	"context"
	"github.com/Clymene-project/Clymene/pkg/multierror"
)

type CompositeWriter struct {
	logWriters []Writer
}

// NewCompositeWriter creates a CompositeWriter
func NewCompositeWriter(logWriters ...Writer) *CompositeWriter {
	return &CompositeWriter{
		logWriters: logWriters,
	}
}

// Writelog calls WriteMetric on each log writer. It will sum up failures, it is not transactional
func (c *CompositeWriter) Writelog(ctx context.Context, tenantID string, batch Batch) (int, int64, int64, error) {
	var errors []error
	for _, writer := range c.logWriters {
		if _, _, _, err := writer.Writelog(ctx, tenantID, batch); err != nil {
			errors = append(errors, err)
		}
	}
	return -1, -1, -1, multierror.Wrap(errors)
}
