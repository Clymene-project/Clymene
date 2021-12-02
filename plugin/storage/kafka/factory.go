package kafka

import (
	"errors"
	"flag"
	"github.com/Clymene-project/Clymene/pkg/kafka/producer"
	"google.golang.org/grpc"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Factory implements storage.Factory and creates write-only storage components backed by kafka.
type Factory struct {
	options Options

	metricsFactory metrics.Factory
	logger         *zap.Logger

	//writeMetrics  WriterMetrics
	//writerSuccess     *app.CountsMetrics
	//writerFail        *app.CountsMetrics
	//writerMarshalFail *app.CountsMetrics

	producer   sarama.AsyncProducer
	marshaller Marshaller
	producer.Builder
}

func (f *Factory) BuildGRPC(conn *grpc.ClientConn) error {
	f.logger.Debug("using kafka")
	return nil
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

	//f.writerSuccess = app.NewKafkaMetricsMgr(metricsFactory, "kafka_written_success")
	//f.writerFail = app.NewKafkaMetricsMgr(metricsFactory, "kafka_written_fail")
	//f.writerMarshalFail = app.NewKafkaMetricsMgr(metricsFactory, "kafka_written_marshal_fail")

	logger.Info("Kafka factory",
		zap.Any("producer builder", f.Builder),
		zap.Any("topic", f.options.TraceTopic))
	p, err := f.NewProducer(f.logger)
	if err != nil {
		return err
	}
	f.producer = p
	switch f.options.Encoding {
	case EncodingProto:
		f.marshaller = newProtobufMarshaller()
	case EncodingJSON:
		f.marshaller = newJSONMarshaller()
	default:
		return errors.New("kafka encoding is not one of '" + EncodingJSON + "' or '" + EncodingProto + "'")
	}
	return nil
}

//func (f *Factory) CreateMetricWriter() (metricstore.Writer, error) {
//	return NewMetricWriter(f.producer, f.marshaller, f.options.MetricTopic, f.writerFail), nil
//}
