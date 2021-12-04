package consumer

import (
	"github.com/Clymene-project/Clymene/cmd/ingester/app/consumer/offset"
	"github.com/Clymene-project/Clymene/cmd/ingester/app/processor"
	"github.com/Clymene-project/Clymene/cmd/ingester/app/processor/decorator"
	"github.com/Clymene-project/Clymene/pkg/kafka/consumer"
	"io"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// ProcessorFactoryParams are the parameters of a ProcessorFactory
type ProcessorFactoryParams struct {
	Parallelism    int
	Topic          string
	BaseProcessor  processor.MetricProcessor
	SaramaConsumer consumer.Consumer
	Factory        metrics.Factory
	Logger         *zap.Logger
	RetryOptions   []decorator.RetryOption
}

// ProcessorFactory is a factory for creating startedProcessors
type ProcessorFactory struct {
	topic          string
	consumer       consumer.Consumer
	metricsFactory metrics.Factory
	logger         *zap.Logger
	baseProcessor  processor.MetricProcessor
	parallelism    int
	retryOptions   []decorator.RetryOption
}

// NewProcessorFactory constructs a new ProcessorFactory
func NewProcessorFactory(params ProcessorFactoryParams) (*ProcessorFactory, error) {
	return &ProcessorFactory{
		topic:          params.Topic,
		consumer:       params.SaramaConsumer,
		metricsFactory: params.Factory,
		logger:         params.Logger,
		baseProcessor:  params.BaseProcessor,
		parallelism:    params.Parallelism,
		retryOptions:   params.RetryOptions,
	}, nil
}

func (c *ProcessorFactory) new(partition int32, minOffset int64) processor.MetricProcessor {
	c.logger.Info("Creating new processors", zap.Int32("partition", partition))

	markOffset := func(offset int64) {
		c.consumer.MarkPartitionOffset(c.topic, partition, offset, "")
	}

	om := offset.NewManager(minOffset, markOffset, partition, c.metricsFactory)

	retryProcessor := decorator.NewRetryingProcessor(c.metricsFactory, c.baseProcessor, c.retryOptions...)
	cp := NewCommittingProcessor(retryProcessor, om)
	metricProcessor := processor.NewDecoratedProcessor(c.metricsFactory, cp)
	pp := processor.NewParallelProcessor(metricProcessor, c.parallelism, c.logger)

	return newStartedProcessor(pp, om)
}

type service interface {
	Start()
	io.Closer
}

type startProcessor interface {
	Start()
	processor.MetricProcessor
}

type startedProcessor struct {
	services  []service
	processor startProcessor
}

func newStartedProcessor(parallelProcessor startProcessor, services ...service) processor.MetricProcessor {
	s := &startedProcessor{
		services:  services,
		processor: parallelProcessor,
	}

	for _, service := range services {
		service.Start()
	}

	s.processor.Start()
	return s
}

func (c *startedProcessor) Process(message processor.Message) error {
	return c.processor.Process(message)
}

func (c *startedProcessor) Close() error {
	c.processor.Close()

	for _, service := range c.services {
		service.Close()
	}
	return nil
}
