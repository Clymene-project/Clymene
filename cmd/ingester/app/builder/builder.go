package builder

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/ingester/app"
	"github.com/Clymene-project/Clymene/cmd/ingester/app/consumer"
	"github.com/Clymene-project/Clymene/cmd/ingester/app/processor"
	kafkaConsumer "github.com/Clymene-project/Clymene/pkg/kafka/consumer"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"strings"
)

func CreateConsumer(logger *zap.Logger, metricsFactory metrics.Factory, metricWriter metricstore.Writer, options app.Options) (*consumer.Consumer, error) {
	var unmarshaller kafka.Unmarshaller
	switch options.Encoding {
	case kafka.EncodingJSON:
		unmarshaller = kafka.NewJSONUnmarshaller()
	case kafka.EncodingProto:
		unmarshaller = kafka.NewProtobufUnmarshaller()
	default:
		return nil, fmt.Errorf(`encoding '%s' not recognised, use one of ("%s")`,
			options.Encoding, strings.Join(kafka.AllEncodings, "\", \""))
	}

	mrParams := &processor.MetricProcessorParams{
		Writer:       metricWriter,
		Unmarshaller: unmarshaller,
	}
	metricProcessor := processor.NewMetricProcessor(mrParams)

	consumerConfig := kafkaConsumer.Configuration{
		Brokers:              options.Brokers,
		Topic:                options.Topic,
		GroupID:              options.GroupID,
		ClientID:             options.ClientID,
		ProtocolVersion:      options.ProtocolVersion,
		AuthenticationConfig: options.AuthenticationConfig,
	}
	saramaConsumer, err := consumerConfig.NewConsumer(logger)
	if err != nil {
		return nil, err
	}

	factoryParams := consumer.ProcessorFactoryParams{
		Topic:          options.Topic,
		Parallelism:    options.Parallelism,
		SaramaConsumer: saramaConsumer,
		BaseProcessor:  metricProcessor,
		Logger:         logger,
		Factory:        metricsFactory,
	}
	processorFactory, err := consumer.NewProcessorFactory(factoryParams)
	if err != nil {
		return nil, err
	}

	consumerParams := consumer.Params{
		InternalConsumer:      saramaConsumer,
		ProcessorFactory:      *processorFactory,
		MetricsFactory:        metricsFactory,
		Logger:                logger,
		DeadlockCheckInterval: options.DeadlockInterval,
	}
	return consumer.New(consumerParams)
}
