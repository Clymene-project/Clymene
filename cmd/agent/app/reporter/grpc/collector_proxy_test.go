// Copyright (c) 2018 The Jaeger Authors.
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

package grpc

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger-lib/metrics/metricstest"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

var _ io.Closer = (*ProxyBuilder)(nil)

func TestMultipleCollectors(t *testing.T) {
	spanHandler1 := &mockSpanHandler{}
	s1, addr1 := initializeGRPCTestServer(t, func(s *grpc.Server) {
		api_v2.RegisterCollectorServiceServer(s, spanHandler1)
	})
	defer s1.Stop()
	spanHandler2 := &mockSpanHandler{}
	s2, addr2 := initializeGRPCTestServer(t, func(s *grpc.Server) {
		api_v2.RegisterCollectorServiceServer(s, spanHandler2)
	})
	defer s2.Stop()

	mFactory := metricstest.NewFactory(time.Microsecond)
	proxy, err := NewCollectorProxy(&ConnBuilder{CollectorHostPorts: []string{addr1.String(), addr2.String()}}, nil, mFactory, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, proxy)
	assert.NotNil(t, proxy.GetReporter())
	assert.NotNil(t, proxy.GetManager())
	assert.NotNil(t, proxy.GetConn())

	var bothServers = false
	r := proxy.GetReporter()
	// TODO do not iterate, just create two batches
	for i := 0; i < 100; i++ {
		err := r.EmitBatch(context.Background(), &jaeger.Batch{Spans: []*jaeger.Span{{OperationName: "op"}}, Process: &jaeger.Process{ServiceName: "service"}})
		require.NoError(t, err)
		if len(spanHandler1.getRequests()) > 0 && len(spanHandler2.getRequests()) > 0 {
			bothServers = true
			break
		}
	}
	c, g := mFactory.Snapshot()
	assert.True(t, len(g) > 0)
	assert.True(t, len(c) > 0)
	assert.Equal(t, true, bothServers)
	require.Nil(t, proxy.Close())
}

func initializeGRPCTestServer(t *testing.T, beforeServe func(server *grpc.Server), opts ...grpc.ServerOption) (*grpc.Server, net.Addr) {
	server := grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	beforeServe(server)
	go func() {
		require.NoError(t, server.Serve(lis))
	}()
	return server, lis.Addr()
}
