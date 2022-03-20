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
	// EncodingJSON is used for timeseries encoded as Protobuf-based JSON.
	EncodingJSON = "json"
	// EncodingProto is used for timeseries encoded as Protobuf.
	EncodingProto = "protobuf"

	configPrefix           = "kafka.producer"
	suffixBrokers          = ".brokers"
	suffixTopic            = ".topic"
	suffixPromtailTopic    = ".promtail.topic"
	suffixEncoding         = ".encoding"
	suffixRequiredAcks     = ".required-acks"
	suffixCompression      = ".compression"
	suffixCompressionLevel = ".compression-level"
	suffixProtocolVersion  = ".protocol-version"
	suffixBatchLinger      = ".batch-linger"
	suffixBatchSize        = ".batch-size"
	suffixBatchMinMessages = ".batch-min-messages"
	suffixBatchMaxMessages = ".batch-max-messages"

	defaultBroker           = "127.0.0.1:9092"
	defaultTopic            = "clymene"
	defaultPromtailTopic    = "clymene-logs"
	defaultEncoding         = EncodingProto
	defaultRequiredAcks     = "local"
	defaultCompression      = "none"
	defaultCompressionLevel = 0
	defaultBatchLinger      = 0
	defaultBatchSize        = 0
	defaultBatchMinMessages = 0
	defaultBatchMaxMessages = 0
)

var (
	// AllEncodings is a list of all supported encodings.
	AllEncodings = []string{EncodingJSON, EncodingProto}

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
	Config        producer.Configuration `mapstructure:",squash"`
	Topic         string                 `mapstructure:"topic"`
	PromtailTopic string                 `mapstructure:"promtail_topic"`
	Encoding      string                 `mapstructure:"encoding"`
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
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
		configPrefix+suffixBatchMinMessages,
		defaultBatchMinMessages,
		"(experimental) The best-effort minimum number of messages needed to send a batch of records to Kafka. Higher value reduce request to Kafka but increase latency and the possibility of data loss in case of process restart. See https://kafka.apache.org/documentation/",
	)
	flagSet.Int(
		configPrefix+suffixBatchMaxMessages,
		defaultBatchMaxMessages,
		"(experimental) Maximum number of message to batch before sending records to Kafka",
	)
	flagSet.String(
		configPrefix+suffixBrokers,
		defaultBroker,
		"The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234'")
	flagSet.String(
		configPrefix+suffixTopic,
		defaultTopic,
		"kafka topic to save time series")
	flagSet.String(
		configPrefix+suffixPromtailTopic,
		defaultPromtailTopic,
		"kafka topic to save log")
	flagSet.String(
		configPrefix+suffixProtocolVersion,
		"",
		"Kafka protocol version - must be supported by kafka server")
	flagSet.String(
		configPrefix+suffixEncoding,
		defaultEncoding,
		fmt.Sprintf(`Encoding of metric ("%s" or "%s") sent to kafka.`, EncodingJSON, EncodingProto),
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
		BatchMinMessages:     v.GetInt(configPrefix + suffixBatchMinMessages),
		BatchMaxMessages:     v.GetInt(configPrefix + suffixBatchMaxMessages),
	}
	opt.Topic = v.GetString(configPrefix + suffixTopic)
	opt.PromtailTopic = v.GetString(configPrefix + suffixPromtailTopic)
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
