package consumer

import (
	"errors"
	"github.com/Clymene-project/Clymene/cmd/ingester/app/processor"
	"io"
)

type comittingProcessor struct {
	processor processor.MetricProcessor
	marker    offsetMarker
	io.Closer
}

type offsetMarker interface {
	MarkOffset(int64)
}

// NewCommittingProcessor returns a processor that commits message offsets to Kafka
func NewCommittingProcessor(processor processor.MetricProcessor, marker offsetMarker) processor.MetricProcessor {
	return &comittingProcessor{
		processor: processor,
		marker:    marker,
	}
}

func (d *comittingProcessor) Process(message processor.Message) error {
	if msg, ok := message.(Message); ok {
		err := d.processor.Process(message)
		if err == nil {
			d.marker.MarkOffset(msg.Offset())
		}
		return err
	}
	return errors.New("committing processor used with non-kafka message")
}
