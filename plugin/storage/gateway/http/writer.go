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
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/pkg/multierror"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

// Writer writes metric to kafka. Implements metricstore.Writer
type Writer struct {
	metrics      WriterMetrics
	logger       *zap.Logger
	client       *http.Client
	url          string
	userAgent    string
	maxErrMsgLen int64
	timeout      time.Duration
	marshaller   Marshaller
}

func (w *Writer) Writelog(ctx context.Context, tenantID string, batch logstore.Batch) (int, int64, int64, error) {
	var bufBytes int64
	var entriesCount64 int64
	var errs []error

	producerMessage := &client.ProducerBatch{TenantID: tenantID, Batch: *batch.(*client.Batch)}
	logsBytes, err := w.marshaller.MarshalLog(producerMessage)
	if err != nil {
		w.metrics.WrittenFailure.Inc(1)
		errs = append(errs, err)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	}

	httpReq, err := http.NewRequest("POST", w.url, bytes.NewReader(logsBytes))
	if err != nil {
		// Errors from NewRequest are from unparsable URLs, so are not
		// recoverable.
		w.logger.Error("NewRequest", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	}
	//httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", w.userAgent)
	//httpReq.Header.Set("Clymene-Version", version.Get().Version)

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	httpResp, err := w.client.Do(httpReq.WithContext(ctx))
	if err != nil {
		// Errors from client.Do are from (for example) network errors, so are
		// recoverable.
		w.logger.Error("client.Do", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	}
	defer func() {
		io.Copy(ioutil.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, w.maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		w.logger.Error("HTTP status error", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return httpResp.StatusCode, bufBytes, entriesCount64, multierror.Wrap(errs)
	}
	if err == nil {
		w.metrics.WrittenSuccess.Inc(1)
	}
	return 201, bufBytes, entriesCount64, multierror.Wrap(errs)

}

func NewLogWriter(
	logger *zap.Logger,
	factory metrics.Factory,
	options Options,
	marshaller Marshaller,
) *Writer {
	writeMetrics := WriterMetrics{
		WrittenSuccess: factory.Counter(metrics.Options{Name: "gateway_http_logs_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: factory.Counter(metrics.Options{Name: "gateway_http_logs_written", Tags: map[string]string{"status": "failure"}}),
	}
	return &Writer{
		metrics:      writeMetrics,
		logger:       logger,
		client:       &http.Client{Transport: newLatencyTransport(http.DefaultTransport, factory), Timeout: options.timeout},
		url:          options.logsUrl,
		userAgent:    options.userAgent,
		maxErrMsgLen: options.maxErrMsgLen,
		timeout:      options.timeout,
		marshaller:   marshaller,
	}
}

func NewMetricWriter(
	logger *zap.Logger,
	factory metrics.Factory,
	options Options,
	marshaller Marshaller,
) *Writer {
	writeMetrics := WriterMetrics{
		WrittenSuccess: factory.Counter(metrics.Options{Name: "gateway_http_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: factory.Counter(metrics.Options{Name: "gateway_http_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	return &Writer{
		metrics:      writeMetrics,
		logger:       logger,
		client:       &http.Client{Transport: newLatencyTransport(http.DefaultTransport, factory), Timeout: options.timeout},
		url:          options.metricsUrl,
		userAgent:    options.userAgent,
		maxErrMsgLen: options.maxErrMsgLen,
		timeout:      options.timeout,
		marshaller:   marshaller,
	}
}

func (w *Writer) WriteMetric(metric []prompb.TimeSeries) error {
	body, err := w.marshaller.MarshalMetric(metric)
	if err != nil {
		w.metrics.WrittenFailure.Inc(1)
		return err
	}
	httpReq, err := http.NewRequest("POST", w.url, bytes.NewReader(body))
	if err != nil {
		// Errors from NewRequest are from unparsable URLs, so are not
		// recoverable.
		w.logger.Error("NewRequest", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return err
	}
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", w.userAgent)
	httpReq.Header.Set("Clymene-Version", version.Get().Version)

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	httpResp, err := w.client.Do(httpReq.WithContext(ctx))
	if err != nil {
		// Errors from client.Do are from (for example) network errors, so are
		// recoverable.
		w.logger.Error("client.Do", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, w.maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		w.logger.Error("HTTP status error", zap.Error(err))
		w.metrics.WrittenFailure.Inc(1)
		return err
	}
	if err == nil {
		w.metrics.WrittenSuccess.Inc(1)
	}
	return err
}

type latencyTransport struct {
	transport http.RoundTripper
	latency   metrics.Timer
	errors    metrics.Counter
}

func (l *latencyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
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
	m := f.Namespace(metrics.NSOptions{Name: "prometheus", Tags: nil})
	return &latencyTransport{
		transport: t,
		latency:   m.Timer(metrics.TimerOptions{Name: "latency", Tags: nil}),
		errors:    m.Counter(metrics.Options{Name: "errors", Tags: nil}),
	}
}
