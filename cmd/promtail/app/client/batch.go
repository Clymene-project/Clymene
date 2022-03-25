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

package client

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	"sort"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
)

type Batch struct {
	Streams   map[string]*logproto.Stream `json:"streams"`
	Bytes     int                         `json:"bytes"`
	CreatedAt time.Time                   `json:"createdAt"`
}

type ProducerBatch struct {
	Batch    Batch  `json:"batch"`
	TenantID string `json:"tenantID"`
}

func newBatch(entries ...api.Entry) *Batch {
	b := &Batch{
		Streams:   map[string]*logproto.Stream{},
		Bytes:     0,
		CreatedAt: time.Now(),
	}

	// Add entries to the Batch
	for _, entry := range entries {
		b.add(entry)
	}

	return b
}

// add an entry to the Batch
func (b *Batch) add(entry api.Entry) {
	b.Bytes += len(entry.Line)

	// Append the entry to an already existing stream (if any)
	labels := labelsMapToString(entry.Labels, ReservedLabelTenantID)
	if stream, ok := b.Streams[labels]; ok {
		stream.Entries = append(stream.Entries, entry.Entry)
		return
	}

	// Add the entry as a new stream
	b.Streams[labels] = &logproto.Stream{
		Labels:  labels,
		Entries: []logproto.Entry{entry.Entry},
	}
}

func labelsMapToString(ls model.LabelSet, without ...model.LabelName) string {
	lstrs := make([]string, 0, len(ls))
Outer:
	for l, v := range ls {
		for _, w := range without {
			if l == w {
				continue Outer
			}
		}
		lstrs = append(lstrs, fmt.Sprintf("%s=%q", l, v))
	}

	sort.Strings(lstrs)
	return fmt.Sprintf("{%s}", strings.Join(lstrs, ", "))
}

// sizeBytes returns the current Batch size in Bytes
func (b *Batch) sizeBytes() int {
	return b.Bytes
}

// sizeBytesAfter returns the size of the Batch after the input entry
// will be added to the Batch itself
func (b *Batch) sizeBytesAfter(entry api.Entry) int {
	return b.Bytes + len(entry.Line)
}

// age of the Batch since its creation
func (b *Batch) age() time.Duration {
	return time.Since(b.CreatedAt)
}

// Encode the Batch as snappy-compressed push request, and returns
// the encoded bytes and the number of encoded entries
func (b *Batch) Encode() ([]byte, int, error) {
	req, entriesCount := b.CreatePushRequest()
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, 0, err
	}
	buf = snappy.Encode(nil, buf)
	return buf, entriesCount, nil
}

// creates push request and returns it, together with number of entries
func (b *Batch) CreatePushRequest() (*logproto.PushRequest, int) {
	req := logproto.PushRequest{
		Streams: make([]logproto.Stream, 0, len(b.Streams)),
	}

	entriesCount := 0
	for _, stream := range b.Streams {
		req.Streams = append(req.Streams, *stream)
		entriesCount += len(stream.Entries)
	}
	return &req, entriesCount
}
