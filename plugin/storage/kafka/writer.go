package kafka

import (
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Shopify/sarama"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type WriterMetrics struct {
	WrittenSuccess metrics.Counter
	WrittenFailure metrics.Counter
}

// Writer writes spans to kafka. Implements spanstore.Writer
type Writer struct {
	//failMetrics *app.CountsMetrics
	producer   sarama.AsyncProducer
	marshaller Marshaller
	topic      string
}

// NewSpanWriter initiates and returns a new kafka spanwriter
func NewSpanWriter(
	producer sarama.AsyncProducer,
	marshaller Marshaller,
	topic string,
	logger *zap.Logger,
	//success *app.CountsMetrics,
	//fail *app.CountsMetrics,
	//marshalFail *app.CountsMetrics,
) *Writer {
	go func() {
		for _ = range producer.Successes() {
			//logger.Debug("kafka writer", zap.String("Topic", msg.Topic))
			//success.CountTopic(msg.Topic, string(msg.Headers[0].Value))
		}
	}()
	go func() {
		for e := range producer.Errors() {
			logger.Error("kafka writer", zap.String("Topic", e.Msg.Topic), zap.Error(e))
			//fail.CountTopic(e.Msg.Topic, string(e.Msg.Headers[0].Value))
		}
	}()

	return &Writer{
		producer:   producer,
		marshaller: marshaller,
		topic:      topic,
		//failMetrics: marshalFail,
	}
}

// Close closes SpanWriter by closing producer
func (w *Writer) Close() error {
	return w.producer.Close()
}

func NewMetricWriter(
	producer sarama.AsyncProducer,
	marshaller Marshaller,
	topic string,
	//fail *app.CountsMetrics,
) *Writer {
	return &Writer{
		producer:   producer,
		marshaller: marshaller,
		topic:      topic,
		//failMetrics: fail,
	}
}

// WriteMetric writes the time series to kafka.
func (w *Writer) WriteMetric(ts []prompb.TimeSeries) error {
	metricsBytes, err := w.marshaller.MarshalMetric(ts)
	if err != nil {
		//w.failMetrics.CountTopic(w.topic, clusterID)
		return err
	}

	// The AsyncProducer accepts messages on a channel and produces them asynchronously
	// in the background as efficiently as possible
	w.producer.Input() <- &sarama.ProducerMessage{
		Topic: w.topic,
		Value: sarama.ByteEncoder(metricsBytes),
	}
	return nil
}
