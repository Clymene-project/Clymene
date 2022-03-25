package builder

import (
	"github.com/Clymene-project/Clymene/cmd/promtail-ingester/app"
	"github.com/Clymene-project/Clymene/cmd/promtail-ingester/app/consumer"
	"github.com/Clymene-project/Clymene/cmd/promtail-ingester/app/processor"
	kafkaConsumer "github.com/Clymene-project/Clymene/pkg/kafka/consumer"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

func CreateConsumer(logger *zap.Logger, metricsFactory metrics.Factory, logWriter logstore.Writer, options app.Options) (*consumer.Consumer, error) {
	var unmarshaller kafka.Unmarshaller
	switch options.Encoding {
	//case kafka.EncodingJSON:
	//	unmarshaller = kafka.NewJSONUnmarshaller()
	//case kafka.EncodingProto:
	//	unmarshaller = kafka.NewProtobufUnmarshaller()
	default:
		unmarshaller = kafka.NewJSONUnmarshaller()
	}

	mrParams := &processor.LogProcessorParams{
		Writer:       logWriter,
		Unmarshaller: unmarshaller,
	}
	logProcessor := processor.NewLogProcessor(mrParams)

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
		BaseProcessor:  logProcessor,
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
