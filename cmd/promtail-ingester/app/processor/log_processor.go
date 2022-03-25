package processor

import (
	"context"
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"io"
)

type LogProcessor interface {
	Process(input Message) error
	io.Closer
}

type Message interface {
	Value() []byte
}

type LogProcessorParams struct {
	Writer       logstore.Writer
	Unmarshaller kafka.Unmarshaller
}

type KafkaLogProcessor struct {
	writer       logstore.Writer
	unmarshaller kafka.Unmarshaller
	io.Closer
}

func NewLogProcessor(params *LogProcessorParams) *KafkaLogProcessor {
	return &KafkaLogProcessor{
		unmarshaller: params.Unmarshaller,
		writer:       params.Writer,
	}
}

// Process unmarshals and writes a single kafka message
func (s *KafkaLogProcessor) Process(message Message) error {
	msg, err := s.unmarshaller.UnmarshalLog(message.Value())
	if err != nil {
		return fmt.Errorf("cannot unmarshall byte array into metrics: %w", err)
	}
	_, _, _, err = s.writer.Writelog(context.Background(), msg.TenantID, &msg.Batch)
	return err
}
