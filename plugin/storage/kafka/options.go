package kafka

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/kafka/auth"
	"github.com/Clymene-project/Clymene/pkg/kafka/producer"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

const (
	// EncodingJSON is used for spans encoded as Protobuf-based JSON.
	EncodingJSON = "json"
	// EncodingProto is used for spans encoded as Protobuf.
	EncodingProto = "protobuf"
	// EncodingZipkinThrift is used for spans encoded as Zipkin Thrift.
	EncodingZipkinThrift = "zipkin-thrift"

	configPrefix  = "kafka.producer"
	suffixBrokers = ".brokers"
	//suffixTopic   = ".topic"

	suffixMetricTopic = "metric.topic"
	suffixTraceTopic  = "trace.topic"

	suffixK8sInfoTopic        = "k8s.info.topic"
	suffixK8sJobInfoTopic     = "k8s.job.topic"
	suffixK8sPodInfoTopic     = "k8s.pod.topic"
	suffixSparseLogTopic      = "k8s.sparse.topic"
	suffixSparseLogModelTopic = "k8s.sparse.model.topic"

	suffixNetstatInfoTopic    = "k8s.netstat.topic"
	suffixPodHistoryInfoTopic = "k8s.pod.history.topic"
	suffixJspdLiteTopic       = "jspd.topic"

	suffixEncoding         = ".encoding"
	suffixRequiredAcks     = ".required-acks"
	suffixCompression      = ".compression"
	suffixCompressionLevel = ".compression-level"
	suffixProtocolVersion  = ".protocol-version"
	suffixBatchLinger      = ".batch-linger"
	suffixBatchSize        = ".batch-size"
	suffixBatchMaxMessages = ".batch-max-messages"

	defaultBroker           = "127.0.0.1:9092"
	traceTopic              = "jaeger-spans"
	metricTopic             = "remote_prom"
	producerMaxMessageBytes = 204857600

	k8sInfoTopic        = "kubernetes_info"
	k8sJobInfoTopic     = "kubernetes_job_info"
	k8sPodInfoTopic     = "kubernetes_pod_info"
	sparseLogTopic      = "sparse_log"
	sparseModelTopic    = "sparse_model"
	netstatInfoTopic    = "netstat_info"
	podHistoryInfoTopic = "pod_history_info" // TODO: to be deleted
	jspdLiteTopic       = "jspd_lite"

	defaultEncoding         = EncodingProto
	defaultRequiredAcks     = "local"
	defaultCompression      = "none"
	defaultCompressionLevel = 0
	defaultBatchLinger      = 0
	defaultBatchSize        = 0
	defaultBatchMaxMessages = 0
)

var (
	// AllEncodings is a list of all supported encodings.
	AllEncodings = []string{EncodingJSON, EncodingProto, EncodingZipkinThrift}

	//requiredAcks is mapping of sarama supported requiredAcks
	requiredAcks = map[string]sarama.RequiredAcks{
		"noack": sarama.NoResponse,
		"local": sarama.WaitForLocal,
		"all":   sarama.WaitForAll,
	}

	// compressionModes is a mapping of supported CompressionType to compressionCodec along with default, min, max compression level
	// https://cwiki.apache.org/confluence/display/KAFKA/KIP-390%3A+Allow+fine-grained+configuration+for+compression
	compressionModes = map[string]struct {
		compressor              sarama.CompressionCodec
		defaultCompressionLevel int
		minCompressionLevel     int
		maxCompressionLevel     int
	}{
		"none": {
			compressor:              sarama.CompressionNone,
			defaultCompressionLevel: 0,
		},
		"gzip": {
			compressor:              sarama.CompressionGZIP,
			defaultCompressionLevel: 6,
			minCompressionLevel:     1,
			maxCompressionLevel:     9,
		},
		"snappy": {
			compressor:              sarama.CompressionSnappy,
			defaultCompressionLevel: 0,
		},
		"lz4": {
			compressor:              sarama.CompressionLZ4,
			defaultCompressionLevel: 9,
			minCompressionLevel:     1,
			maxCompressionLevel:     17,
		},
		"zstd": {
			compressor:              sarama.CompressionZSTD,
			defaultCompressionLevel: 3,
			minCompressionLevel:     -131072,
			maxCompressionLevel:     22,
		},
	}
)

// Options stores the configuration options for Kafka
type Options struct {
	Config      producer.Configuration `mapstructure:",squash"`
	TraceTopic  string                 `mapstructure:"topic"`
	MetricTopic string                 `mapstructure:"metrictopic"`

	K8sInfoTopic        string
	K8sJobInfoTopic     string
	K8sPodInfoTopic     string
	SparseLogTopic      string
	SparseModelTopic    string
	NetstatInfoTopic    string
	PodHistoryInfoTopic string // TODO: to be deleted

	JspdTopic string

	Encoding string `mapstructure:"encoding"`
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixBrokers,
		defaultBroker,
		"The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234'")
	//flagSet.String(
	//	configPrefix+suffixTopic,
	//	traceTopic,
	//	"The name of the kafka topic")
	flagSet.String(
		configPrefix+suffixProtocolVersion,
		"",
		"Kafka protocol version - must be supported by kafka server")
	flagSet.String(
		configPrefix+suffixEncoding,
		defaultEncoding,
		fmt.Sprintf(`Encoding of spans ("%s" or "%s") sent to kafka.`, EncodingJSON, EncodingProto),
	)
	flagSet.String(
		configPrefix+suffixRequiredAcks,
		defaultRequiredAcks,
		"(experimental) Required kafka broker acknowledgement. i.e. noack, local, all",
	)
	flagSet.String(
		configPrefix+suffixCompression,
		defaultCompression,
		"(experimental) Type of compression (none, gzip, snappy, lz4, zstd) to use on messages",
	)
	flagSet.String(
		configPrefix+suffixMetricTopic,
		metricTopic,
		"metric topic change option(default = "+metricTopic+")",
	)
	flagSet.String(
		configPrefix+suffixK8sInfoTopic,
		k8sInfoTopic,
		"kubernetes_info topic change option(default = "+k8sInfoTopic+")",
	)
	flagSet.String(
		configPrefix+suffixK8sJobInfoTopic,
		k8sJobInfoTopic,
		"kubernetes_job_info topic change option(default = "+k8sJobInfoTopic+")",
	)
	flagSet.String(
		configPrefix+suffixK8sPodInfoTopic,
		k8sPodInfoTopic,
		"kubernetes_pod_info topic change option(default = "+k8sPodInfoTopic+")",
	)
	flagSet.String(
		configPrefix+suffixSparseLogTopic,
		sparseLogTopic,
		"sparse_log topic change option(default = "+sparseLogTopic+")",
	)

	flagSet.String(
		configPrefix+suffixSparseLogModelTopic,
		sparseModelTopic,
		"sparse_model topic change option(default = "+sparseModelTopic+")",
	)

	flagSet.String(
		configPrefix+suffixNetstatInfoTopic,
		netstatInfoTopic,
		"netstat_info topic change option(default = "+netstatInfoTopic+")",
	)
	flagSet.String(
		configPrefix+suffixPodHistoryInfoTopic,
		podHistoryInfoTopic,
		"pod_history_info topic change option(default = "+podHistoryInfoTopic+")",
	)
	flagSet.String(
		configPrefix+suffixJspdLiteTopic,
		jspdLiteTopic,
		"jspd_lite topic change option(default = "+jspdLiteTopic+")",
	)

	flagSet.String(
		configPrefix+suffixTraceTopic,
		traceTopic,
		"trace topic change option(default = "+traceTopic+")",
	)

	flagSet.Int(
		configPrefix+suffixCompressionLevel,
		defaultCompressionLevel,
		"(experimental) compression level to use on messages. gzip = 1-9 (default = 6), snappy = none, lz4 = 1-17 (default = 9), zstd = -131072 - 22 (default = 3)",
	)
	flagSet.Duration(
		configPrefix+suffixBatchLinger,
		defaultBatchLinger,
		"(experimental) Time interval to wait before sending records to Kafka. Higher value reduce request to Kafka but increase latency and the possibility of data loss in case of process restart. See https://kafka.apache.org/documentation/",
	)
	flagSet.Int(
		configPrefix+suffixBatchSize,
		defaultBatchSize,
		"(experimental) Number of bytes to batch before sending records to Kafka. Higher value reduce request to Kafka but increase latency and the possibility of data loss in case of process restart. See https://kafka.apache.org/documentation/",
	)
	flagSet.Int(
		configPrefix+suffixBatchMaxMessages,
		defaultBatchMaxMessages,
		"(experimental) Number of message to batch before sending records to Kafka. Higher value reduce request to Kafka but increase latency and the possibility of data loss in case of process restart. See https://kafka.apache.org/documentation/",
	)
	auth.AddFlags(configPrefix, flagSet)
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	authenticationOptions := auth.AuthenticationConfig{}
	authenticationOptions.InitFromViper(configPrefix, v)

	requiredAcks, err := getRequiredAcks(v.GetString(configPrefix + suffixRequiredAcks))
	if err != nil {
		log.Fatal(err)
	}

	compressionMode := strings.ToLower(v.GetString(configPrefix + suffixCompression))
	compressionModeCodec, err := getCompressionMode(compressionMode)
	if err != nil {
		log.Fatal(err)
	}

	compressionLevel, err := getCompressionLevel(compressionMode, v.GetInt(configPrefix+suffixCompressionLevel))
	if err != nil {
		log.Fatal(err)
	}

	opt.Config = producer.Configuration{
		Brokers:              strings.Split(stripWhiteSpace(v.GetString(configPrefix+suffixBrokers)), ","),
		RequiredAcks:         requiredAcks,
		Compression:          compressionModeCodec,
		CompressionLevel:     compressionLevel,
		ProtocolVersion:      v.GetString(configPrefix + suffixProtocolVersion),
		AuthenticationConfig: authenticationOptions,
		BatchLinger:          v.GetDuration(configPrefix + suffixBatchLinger),
		BatchSize:            v.GetInt(configPrefix + suffixBatchSize),
		BatchMaxMessages:     v.GetInt(configPrefix + suffixBatchMaxMessages),
		MaxMessageBytes:      producerMaxMessageBytes,
	}

	opt.TraceTopic = v.GetString(configPrefix + suffixTraceTopic)
	opt.MetricTopic = v.GetString(configPrefix + suffixMetricTopic)

	opt.K8sInfoTopic = v.GetString(configPrefix + suffixK8sInfoTopic)
	opt.K8sJobInfoTopic = v.GetString(configPrefix + suffixK8sJobInfoTopic)
	opt.K8sPodInfoTopic = v.GetString(configPrefix + suffixK8sPodInfoTopic)
	opt.SparseLogTopic = v.GetString(configPrefix + suffixSparseLogTopic)
	opt.SparseModelTopic = v.GetString(configPrefix + suffixSparseLogModelTopic)
	opt.NetstatInfoTopic = v.GetString(configPrefix + suffixNetstatInfoTopic)
	opt.PodHistoryInfoTopic = v.GetString(configPrefix + suffixPodHistoryInfoTopic)

	opt.JspdTopic = v.GetString(configPrefix + suffixJspdLiteTopic)

	opt.Encoding = v.GetString(configPrefix + suffixEncoding)
}

// stripWhiteSpace removes all whitespace characters from a string
func stripWhiteSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// getCompressionLevel to get compression level from compression type
func getCompressionLevel(mode string, compressionLevel int) (int, error) {
	compressionModeData, ok := compressionModes[mode]
	if !ok {
		return 0, fmt.Errorf("cannot find compression mode for compressionMode %v", mode)
	}

	if compressionLevel == defaultCompressionLevel {
		return compressionModeData.defaultCompressionLevel, nil
	}

	if compressionModeData.minCompressionLevel > compressionLevel || compressionModeData.maxCompressionLevel < compressionLevel {
		return 0, fmt.Errorf("compression level %d for '%s' is not within valid range [%d, %d]", compressionLevel, mode, compressionModeData.minCompressionLevel, compressionModeData.maxCompressionLevel)
	}

	return compressionLevel, nil
}

//getCompressionMode maps input modes to sarama CompressionCodec
func getCompressionMode(mode string) (sarama.CompressionCodec, error) {
	compressionMode, ok := compressionModes[mode]
	if !ok {
		return 0, fmt.Errorf("unknown compression mode: %v", mode)
	}

	return compressionMode.compressor, nil
}

//getRequiredAcks maps input ack values to sarama requiredAcks
func getRequiredAcks(acks string) (sarama.RequiredAcks, error) {
	requiredAcks, ok := requiredAcks[strings.ToLower(acks)]
	if !ok {
		return 0, fmt.Errorf("unknown Required Ack: %s", acks)
	}
	return requiredAcks, nil
}
