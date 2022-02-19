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
)

type Writer interface {
	//WriterMetric
	Writelog(ctx context.Context, tenantID string, batch []byte) (int, error)
}

type WriterMetric interface {
	EncodedBytesInc(int64)
	SentBytesInc(int64)
	SentEntriesInc(int64)
	StreamLagSet(int64)
	StreamLagInit()
	RequestDurationSet(float64)
	BatchRetriesInc()
	DroppedBytesInc(int64)
	DroppedEntriesInc(int64)
}

type Batch interface {
	Encode() ([]byte, int, error)
}
