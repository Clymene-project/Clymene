package client

import (
	"context"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/agent/app/parser"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"

	"github.com/grafana/dskit/backoff"
	"github.com/prometheus/common/model"
)

const (
	contentType  = "application/x-protobuf"
	maxErrMsgLen = 1024

	// Label reserved to override the tenant ID while processing
	// pipeline stages
	ReservedLabelTenantID = "__tenant_id__"

	LatencyLabel = "filename"
	HostLabel    = "host"
)

var UserAgent = fmt.Sprintf("promtail/%s", version.Get().Version)

// Client pushes entries to Loki and can be stopped
type Client interface {
	api.EntryHandler
	// StopNow Stop goroutine sending batch of entries without retries.
	StopNow()
}

// Client for pushing logs in snappy-compressed protos over HTTP.
type client struct {
	logger *zap.Logger

	entries chan api.Entry

	once sync.Once
	wg   sync.WaitGroup

	// ctx is used in any upstream calls from the `client`.
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter logstore.Writer
	options   Options

	externalLabels model.LabelSet
	writerMetrics  WriterMetrics
}

// Tripperware can wrap a roundtripper.
type Tripperware func(http.RoundTripper) http.RoundTripper

type WriterMetrics struct {
	EncodedBytes    metrics.Counter
	SentBytes       metrics.Counter
	DroppedBytes    metrics.Counter
	SentEntries     metrics.Counter
	DroppedEntries  metrics.Counter
	RequestDuration metrics.Histogram
	BatchRetries    metrics.Counter
	StreamLag       metrics.Gauge
}

func NewWriterMetrics(metricFactory metrics.Factory) WriterMetrics {
	return WriterMetrics{
		EncodedBytes: metricFactory.Counter(
			metrics.Options{
				Name: "encoded_bytes_total",
				Help: "Number of bytes encoded and ready to send.",
			}),
		SentBytes: metricFactory.Counter(
			metrics.Options{
				Name: "sent_bytes_total",
				Help: "Number of bytes sent.",
			}),
		DroppedBytes: metricFactory.Counter(
			metrics.Options{
				Name: "dropped_bytes_total",
				Help: "Number of bytes dropped because failed to be sent to the ingester after all retries.",
			}),
		SentEntries: metricFactory.Counter(
			metrics.Options{
				Name: "sent_entries_total",
				Help: "Number of log entries sent to the ingester.",
			}),
		DroppedEntries: metricFactory.Counter(
			metrics.Options{
				Name: "dropped_entries_total",
				Help: "Number of log entries dropped because failed to be sent to the ingester after all retries.",
			}),
		RequestDuration: metricFactory.Histogram(
			metrics.HistogramOptions{
				Name: "request_duration_seconds",
				Help: "Duration of send requests.",
				Tags: map[string]string{"status_code": HostLabel}}),
		BatchRetries: metricFactory.Counter(
			metrics.Options{
				Name: "batch_retries_total",
				Help: "Number of times batches has had to be retried.",
			}),
		StreamLag: metricFactory.Gauge(
			metrics.Options{
				Name: "stream_lag_seconds",
				Help: "Difference between current time and last batch timestamp for successful sends",
			}),
	}
}

// New makes a new Client.
func New(options Options, logWriter logstore.Writer, factory metrics.Factory, logger *zap.Logger) (Client, error) {
	return newClient(options, logWriter, factory, logger)
}

func newClient(options Options, logWriter logstore.Writer, metricFactory metrics.Factory, logger *zap.Logger) (*client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &client{
		logger:         logger.With(zap.String("component", "client")),
		entries:        make(chan api.Entry),
		logWriter:      logWriter,
		ctx:            ctx,
		cancel:         cancel,
		options:        options,
		externalLabels: options.ExternalLabels.LabelSet,
		writerMetrics:  NewWriterMetrics(metricFactory),
	}
	c.wg.Add(1)
	go c.run()
	return c, nil
}

func (c *client) run() {
	batches := map[string]*batch{}

	// Given the client handles multiple batches (1 per tenant) and each batch
	// can be created at a different point in time, we look for batches whose
	// max wait time has been reached every 10 times per BatchWait, so that the
	// maximum delay we have sending batches is 10% of the max waiting time.
	// We apply a cap of 10ms to the ticker, to avoid too frequent checks in
	// case the BatchWait is very low.
	minWaitCheckFrequency := 10 * time.Millisecond
	maxWaitCheckFrequency := c.options.BatchWait / 10
	if maxWaitCheckFrequency < minWaitCheckFrequency {
		maxWaitCheckFrequency = minWaitCheckFrequency
	}

	maxWaitCheck := time.NewTicker(maxWaitCheckFrequency)
	defer func() {
		maxWaitCheck.Stop()
		// Send all pending batches
		for tenantID, batch := range batches {
			c.sendBatch(tenantID, batch)
		}

		c.wg.Done()
	}()

	for {
		select {
		case e, ok := <-c.entries:
			if !ok {
				return
			}
			e, tenantID := c.processEntry(e)
			batch, ok := batches[tenantID]

			// If the batch doesn't exist yet, we create a new one with the entry
			if !ok {
				batches[tenantID] = newBatch(e)
				break
			}

			// If adding the entry to the batch will increase the size over the max
			// size allowed, we do send the current batch and then create a new one

			if batch.sizeBytesAfter(e) > c.options.BatchSize {
				c.sendBatch(tenantID, batch)

				batches[tenantID] = newBatch(e)
				break
			}

			// The max size of the batch isn't reached, so we can add the entry
			batch.add(e)

		case <-maxWaitCheck.C:
			// Send all batches whose max wait time has been reached
			for tenantID, batch := range batches {
				if batch.age() < c.options.BatchWait {
					continue
				}

				c.sendBatch(tenantID, batch)
				delete(batches, tenantID)
			}
		}
	}
}

func (c *client) Chan() chan<- api.Entry {
	return c.entries
}

func (c *client) sendBatch(tenantID string, batch *batch) {
	buf, entriesCount, err := batch.Encode()
	if err != nil {
		c.logger.Error("error encoding batch", zap.Error(err))
		return
	}
	bufBytes := int64(len(buf))
	c.writerMetrics.EncodedBytes.Inc(bufBytes)

	backoff := backoff.New(c.ctx, c.options.BackoffConfig)
	var status int
	for {
		start := time.Now()
		// send uses `timeout` internally, so `context.Background` is good enough.
		status, err := c.logWriter.Writelog(context.Background(), tenantID, buf)

		c.writerMetrics.RequestDuration.Record(time.Since(start).Seconds())
		if err == nil {
			c.writerMetrics.SentBytes.Inc(bufBytes)
			c.writerMetrics.SentEntries.Inc(int64(entriesCount))

			for _, s := range batch.streams {
				lbls, err := parser.ParseMetric(s.Labels)
				if err != nil {
					// is this possible?
					c.logger.Warn("error converting stream label string to label.Labels, cannot update lagging metric", zap.Error(err))
					return
				}
				var lblSet model.LabelSet
				for i := range lbls {
					for _, lbl := range c.options.StreamLagLabels {
						if lbls[i].Name == lbl {
							if lblSet == nil {
								lblSet = model.LabelSet{}
							}

							lblSet = lblSet.Merge(model.LabelSet{
								model.LabelName(lbl): model.LabelValue(lbls[i].Value),
							})
						}
					}
				}
				if lblSet != nil {
					c.writerMetrics.StreamLag.Update(int64(time.Since(s.Entries[len(s.Entries)-1].Timestamp).Seconds()))
				}
			}
			return
		}

		// Only retry 429s, 500s and connection-level errors.
		if status > 0 && status != 429 && status/100 != 5 {
			break
		}

		c.logger.Warn("error sending batch, will retry", zap.Int("status", status), zap.Error(err))
		c.writerMetrics.BatchRetries.Inc(1)
		backoff.Wait()

		// Make sure it sends at least once before checking for retry.
		if !backoff.Ongoing() {
			break
		}
	}

	if err != nil {
		c.logger.Error("final error sending batch", zap.Int("status", status), zap.Error(err))
		c.writerMetrics.DroppedBytes.Inc(bufBytes)
		c.writerMetrics.DroppedEntries.Inc(int64(entriesCount))
	}
}
func (c *client) getTenantID(labels model.LabelSet) string {
	// Check if it has been overridden while processing the pipeline stages
	if value, ok := labels[ReservedLabelTenantID]; ok {
		return string(value)
	}

	// Check if has been specified in the config
	if c.options.TenantID != "" {
		return c.options.TenantID
	}

	// Defaults to an empty string, which means the X-Scope-OrgID header
	// will not be sent
	return ""
}

// Stop the client.
func (c *client) Stop() {
	c.once.Do(func() { close(c.entries) })
	c.wg.Wait()
}

// StopNow stops the client without retries
func (c *client) StopNow() {
	// cancel will stop retrying http requests.
	c.cancel()
	c.Stop()
}

func (c *client) processEntry(e api.Entry) (api.Entry, string) {
	if len(c.externalLabels) > 0 {
		e.Labels = c.externalLabels.Merge(e.Labels)
	}
	tenantID := c.getTenantID(e.Labels)
	return e, tenantID
}

func (c *client) UnregisterLatencyMetric(labels model.LabelSet) {
	c.writerMetrics.StreamLag.Update(0)
}
