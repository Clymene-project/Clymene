package kafka

import (
	"errors"
	"flag"
	"github.com/Clymene-project/Clymene/pkg/kafka/producer"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
)

// Factory implements storage.Factory and creates write-only storage components backed by kafka.
type Factory struct {
	options Options

	metricsFactory metrics.Factory
	logger         *zap.Logger

	producer   sarama.AsyncProducer
	marshaller Marshaller
	producer.Builder
}

func (f *Factory) CreateLogWriter() (logstore.Writer, error) {
	return NewLogWriter(f.producer, NewJSONMarshaller(), f.options.PromtailTopic, f.metricsFactory, f.logger), nil
}

// NewFactory creates a new Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.options.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.options.InitFromViper(v)
	f.Builder = &f.options.Config
}

// InitFromOptions initializes factory from options.
func (f *Factory) InitFromOptions(o Options) {
	f.options = o
	f.Builder = &f.options.Config
}

// Initialize implements storage.Factory
func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	logger.Info("Factory Initialize", zap.String("type", "kafka"))

	logger.Info("Kafka factory",
		zap.Any("producer builder", f.Builder),
		zap.Any("topic", f.options.Topic),
		zap.Any("topic of promtail", f.options.PromtailTopic))

	p, err := f.NewProducer(logger)
	if err != nil {
		return err
	}
	f.producer = p
	switch f.options.Encoding {
	case EncodingProto:
		f.marshaller = NewProtobufMarshaller()
		logger.Info("promtail can only use json Marshaller.")
	case EncodingJSON:
		if f.options.Flatten {
			f.marshaller = NewJsonFlattenMarshaller()
		} else {
			f.marshaller = NewJSONMarshaller()
		}
	default:
		return errors.New("kafka encoding is not one of '" + EncodingJSON + "' or '" + EncodingProto + "'")
	}
	return nil
}

func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
	return NewMetricWriter(f.producer, f.marshaller, f.options.Topic, f.metricsFactory, f.logger), nil
}

var _ io.Closer = (*Factory)(nil)

// Close closes the resources held by the factory
func (f *Factory) Close() error {
	return f.options.Config.TLS.Close()
}
