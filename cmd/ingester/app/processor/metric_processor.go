package processor

import (
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"io"
)

type MetricProcessor interface {
	Process(input Message) error
	io.Closer
}

type Message interface {
	Value() []byte
}

type MetricProcessorParams struct {
	Writer       metricstore.Writer
	Unmarshaller kafka.Unmarshaller
}

type KafkaMetricProcessor struct {
	writer       metricstore.Writer
	unmarshaller kafka.Unmarshaller
	io.Closer
}

func NewMetricProcessor(params *MetricProcessorParams) *KafkaMetricProcessor {
	return &KafkaMetricProcessor{
		unmarshaller: params.Unmarshaller,
		writer:       params.Writer,
	}
}

// Process unmarshals and writes a single kafka message
func (s *KafkaMetricProcessor) Process(message Message) error {
	metrics, err := s.unmarshaller.Unmarshal(message.Value())
	if err != nil {
		return fmt.Errorf("cannot unmarshall byte array into metrics: %w", err)
	}
	return s.writer.WriteMetric(metrics)
}
