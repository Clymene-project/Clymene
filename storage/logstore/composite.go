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

// Writelog calls WriteMetric on each metric writer. It will sum up failures, it is not transactional
func (c *CompositeWriter) Writelog() error {
	var errors []error
	for _, writer := range c.logWriters {
		if err := writer.Writelog(); err != nil {
			errors = append(errors, err)
		}
	}
	return multierror.Wrap(errors)
}
