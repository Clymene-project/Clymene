package kafka

import (
	"context"
	"github.com/Clymene-project/Clymene/pkg/multierror"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Shopify/sarama"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

// Writer writes metric to kafka. Implements metricstore.Writer
type Writer struct {
	metrics    WriterMetrics
	producer   sarama.AsyncProducer
	marshaller Marshaller
	topic      string
}

func (w *Writer) Writelog(ctx context.Context, tenantID string, batch logstore.Batch) (int, int64, int64, error) {
	var bufBytes int64
	var entriesCount64 int64
	var errs []error

	req, entriesCount := batch.CreatePushRequest()
	entriesCount64 = int64(entriesCount)
	bufBytes = int64(len(req.Streams)) // Put the length value of Streams for data verification

	metricsBytes, err := w.marshaller.MarshalLog(req)
	if err != nil {
		w.metrics.WrittenFailure.Inc(1)
		errs = append(errs, err)
		return 201, bufBytes, entriesCount64, multierror.Wrap(errs)
	}

	// The AsyncProducer accepts messages on a channel and produces them asynchronously
	// in the background as efficiently as possible
	// If there is no key provided, then Kafka will partition the data in a round-robin fashion.
	w.producer.Input() <- &sarama.ProducerMessage{
		Topic:    w.topic,
		Metadata: tenantID,
		Value:    sarama.ByteEncoder(metricsBytes),
	}
	return 201, bufBytes, entriesCount64, multierror.Wrap(errs)

}

// Close closes metricWriter by closing producer
func (w *Writer) Close() error {
	return w.producer.Close()
}

func NewMetricWriter(
	producer sarama.AsyncProducer,
	marshaller Marshaller,
	topic string,
	factory metrics.Factory,
	logger *zap.Logger,
) *Writer {
	writeMetrics := WriterMetrics{
		WrittenSuccess: factory.Counter(metrics.Options{Name: "kafka_metrics_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: factory.Counter(metrics.Options{Name: "kafka_metrics_written", Tags: map[string]string{"status": "failure"}}),
	}
	go func() {
		for range producer.Successes() {
			writeMetrics.WrittenSuccess.Inc(1)
		}
	}()
	go func() {
		for e := range producer.Errors() {
			if e != nil && e.Err != nil {
				logger.Error(e.Err.Error())
			}
			writeMetrics.WrittenFailure.Inc(1)
		}
	}()

	return &Writer{
		producer:   producer,
		marshaller: marshaller,
		topic:      topic,
	}
}

// WriteMetric writes the time series to kafka.
func (w *Writer) WriteMetric(ts []prompb.TimeSeries) error {
	metricsBytes, err := w.marshaller.MarshalMetric(ts)
	if err != nil {
		w.metrics.WrittenFailure.Inc(1)
		return err
	}

	// The AsyncProducer accepts messages on a channel and produces them asynchronously
	// in the background as efficiently as possible
	// If there is no key provided, then Kafka will partition the data in a round-robin fashion.
	w.producer.Input() <- &sarama.ProducerMessage{
		Topic: w.topic,
		Value: sarama.ByteEncoder(metricsBytes),
	}
	return nil
}

func NewLogWriter(
	producer sarama.AsyncProducer,
	marshaller Marshaller,
	topic string,
	factory metrics.Factory,
	logger *zap.Logger,
) *Writer {
	writeLogs := WriterMetrics{
		WrittenSuccess: factory.Counter(metrics.Options{Name: "kafka_logs_written", Tags: map[string]string{"status": "success"}}),
		WrittenFailure: factory.Counter(metrics.Options{Name: "kafka_logs_written", Tags: map[string]string{"status": "failure"}}),
	}
	go func() {
		for range producer.Successes() {
			writeLogs.WrittenSuccess.Inc(1)
		}
	}()
	go func() {
		for e := range producer.Errors() {
			if e != nil && e.Err != nil {
				logger.Error(e.Err.Error())
			}
			writeLogs.WrittenFailure.Inc(1)
		}
	}()

	return &Writer{
		producer:   producer,
		marshaller: marshaller,
		topic:      topic,
	}
}
