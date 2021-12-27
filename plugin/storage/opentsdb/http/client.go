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

package http

import (
	"bufio"
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	hosts         []string
	authorization string
	l             *zap.Logger
	c             *http.Client
	timeout       time.Duration
	converter     *dbmodel.Converter
	writeMetrics  WriterMetrics
	maxChunk      int
}

type Options struct {
	Hosts        []Hosts
	HttpPassword string
	HttpUsername string
	SSL          bool
	HttpApiPath  string
	Factory      metrics.Factory
	Timeout      time.Duration
	MaxChunk     int
}

type Hosts struct {
	Host string
	Port int
}

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

// http://opentsdb.net/docs/build/html/user_guide/configuration.html

func (c *Client) SendData(metrics []prompb.TimeSeries) error {
	// tsd.http.request.max_chunk
	// The maximum request body size to support for incoming HTTP requests when chunking is enabled. 4096
	// TODO Sent normally from max chunk value of 50 ...
	q := len(metrics) / c.maxChunk
	r := len(metrics) % c.maxChunk
	if r != 0 {
		q += 1
	}
	for i := 1; i <= q; i++ {
		var timeSeriesDiv []prompb.TimeSeries
		if i == 1 {
			timeSeriesDiv = metrics[:i*c.maxChunk]
		} else if i != q {
			timeSeriesDiv = metrics[(i-1)*c.maxChunk : i*c.maxChunk]
		} else {
			timeSeriesDiv = metrics[(i-1)*c.maxChunk:]
		}
		jsonTS, err := c.converter.ConvertTsToOpenTSDBJSON(timeSeriesDiv)
		if err != nil {
			c.writeMetrics.WrittenFailure.Inc(int64(len(c.hosts)))
			return err
		}
		for _, host := range c.hosts {
			c.makeRequest(host, jsonTS)
		}
	}

	return nil
}
func (c *Client) makeRequest(host string, json []byte) {
	var MaxErrMsgLen int64
	MaxErrMsgLen = 256
	httpReq, err := http.NewRequest("POST", host, bytes.NewBuffer(json))
	if err != nil {
		c.writeMetrics.WrittenFailure.Inc(1)
		c.l.Error("NewRequest", zap.Error(err))
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.authorization != "" {
		httpReq.Header.Set("Authorization", c.authorization)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	httpResp, err := c.c.Do(httpReq.WithContext(ctx))
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.l.Warn("EOF Error", zap.String("max-chunk", "Adjust the maxChunk option value"))
		}
		c.writeMetrics.WrittenFailure.Inc(1)
		c.l.Error("client.Do", zap.Error(err))
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, MaxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		c.writeMetrics.WrittenFailure.Inc(1)
		c.l.Error("HTTP status error", zap.Error(err))
		return
	}
	if err == nil {
		c.writeMetrics.WrittenSuccess.Inc(1)
	} else {
		c.l.Error("HTTP status error", zap.Error(err))
	}
}

func NewClient(o *Options, converter *dbmodel.Converter, l *zap.Logger) *Client {
	protocol := "http"
	if o.SSL {
		protocol = "https"
	}
	var hosts []string
	for _, h := range o.Hosts {
		hosts = append(hosts, fmt.Sprintf("%s://%s:%d/%s", protocol, h.Host, h.Port, o.HttpApiPath))
	}
	authorization := ""
	if o.HttpUsername != "" && o.HttpPassword != "" {
		encodeData := o.HttpUsername + ":" + o.HttpPassword
		authorization = "Basic " + b64.StdEncoding.EncodeToString([]byte(encodeData))
	}
	writeMetrics := WriterMetrics{
		WrittenSuccess: o.Factory.Counter(metrics.Options{Name: "opentsdb_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: o.Factory.Counter(metrics.Options{Name: "opentsdb_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	return &Client{
		hosts:         hosts,
		authorization: authorization,
		c:             &http.Client{Transport: newLatencyTransport(http.DefaultTransport, o.Factory), Timeout: o.Timeout},
		converter:     converter,
		writeMetrics:  writeMetrics,
		timeout:       o.Timeout,
		l:             l,
		maxChunk:      o.MaxChunk,
	}
}

type latencyTransport struct {
	transport http.RoundTripper
	latency   metrics.Timer
	errors    metrics.Counter
}

func (l latencyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	now := time.Now()
	resp, err := l.transport.RoundTrip(request)
	if err != nil {
		l.errors.Inc(1)
		return resp, err
	}
	l.latency.Record(time.Since(now))
	return resp, err
}

func newLatencyTransport(t http.RoundTripper, f metrics.Factory) http.RoundTripper {
	m := f.Namespace(metrics.NSOptions{Name: "opentsdb", Tags: nil})
	return &latencyTransport{
		transport: t,
		latency:   m.Timer(metrics.TimerOptions{Name: "latency", Tags: nil}),
		errors:    m.Counter(metrics.Options{Name: "errors", Tags: nil}),
	}
}
