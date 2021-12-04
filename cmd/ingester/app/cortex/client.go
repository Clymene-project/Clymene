package cortex

import (
	"bufio"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Writer struct {
	url          string
	userAgent    string
	client       *http.Client
	timeout      time.Duration
	l            *zap.Logger
	maxErrMsgLen int64
}

func NewWriter(logger *zap.Logger, f metrics.Factory) *Writer {
	latencyTransport := newLatencyTransport(http.DefaultTransport, f)
	cortexDistributor := strings.ReplaceAll(os.Getenv("HTTP_PUSH"), " ", "")
	if cortexDistributor == "" {
		logger.Fatal("cortex connection err: check cortex distributor url")
	}
	return &Writer{
		url:          os.Getenv("HTTP_PUSH"),
		client:       &http.Client{Transport: latencyTransport},
		l:            logger,
		timeout:      30 * time.Second,
		userAgent:    "Prometheus/2.17.2", // cortex 수신시 버전체크용 하드코딩
		maxErrMsgLen: 256,
	}
}

type latencyTransport struct {
	transport http.RoundTripper
	latency   metrics.Timer
	errors    metrics.Counter
}

// RoundTrip response metric 을 남기기 위한 인터페이스 구현
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
	m := f.Namespace(metrics.NSOptions{Name: "cortex-client", Tags: nil})
	return &latencyTransport{
		transport: t,
		latency:   m.Timer(metrics.TimerOptions{Name: "latency", Tags: nil}),
		errors:    m.Counter(metrics.Options{Name: "errors", Tags: nil}),
	}
}

func (c *Writer) WriteMetric(Metrics []byte) error {
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewReader(Metrics))
	if err != nil {
		// Errors from NewRequest are from unparsable URLs, so are not
		// recoverable.
		c.l.Error("NewRequest", zap.Error(err))
		return err
	}
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	httpResp, err := c.client.Do(httpReq.WithContext(ctx))
	if err != nil {
		// Errors from client.Do are from (for example) network errors, so are
		// recoverable.
		c.l.Error("client.Do", zap.Error(err))
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, c.maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		c.l.Error("HTTP status error", zap.Error(err))
		return err
	}
	return err
}
