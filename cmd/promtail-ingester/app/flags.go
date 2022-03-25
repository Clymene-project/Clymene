package app

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/kafka/auth"
	kafkaConsumer "github.com/Clymene-project/Clymene/pkg/kafka/consumer"
	"github.com/Clymene-project/Clymene/plugin/storage/kafka"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// ConfigPrefix is a prefix for the ingester flags
	ConfigPrefix = "clymene-ingester"
	// KafkaConsumerConfigPrefix is a prefix for the Kafka flags
	KafkaConsumerConfigPrefix = "kafka.consumer"
	// SuffixBrokers is a suffix for the brokers flag
	SuffixBrokers = ".brokers"
	// SuffixPromtailTopic is a suffix for the promtail topic flag
	SuffixPromtailTopic = ".promtail.topic"
	// SuffixGroupID is a suffix for the group-id flag
	SuffixGroupID = ".group-id"
	// SuffixClientID is a suffix for the client-id flag
	SuffixClientID = ".client-id"
	// SuffixProtocolVersion Kafka protocol version - must be supported by kafka server
	SuffixProtocolVersion = ".protocol-version"
	// SuffixEncoding is a suffix for the encoding flag
	SuffixEncoding = ".encoding"
	// SuffixDeadlockInterval is a suffix for deadlock detecor flag
	SuffixDeadlockInterval = ".deadlockInterval"
	// SuffixParallelism is a suffix for the parallelism flag
	SuffixParallelism = ".parallelism"
	// SuffixHTTPPort is a suffix for the HTTP port
	SuffixHTTPPort = ".http-port"
	// DefaultBroker is the default kafka broker
	DefaultBroker = "127.0.0.1:9092"
	// DefaultTopic is the default kafka topic
	DefaultTopic = "clymene"
	// DefaultPromtailTopic is the default kafka topic for promtail
	DefaultPromtailTopic = "clymene-logs"
	// DefaultGroupID is the default consumer Group ID
	DefaultGroupID = "clymene"
	// DefaultClientID is the default consumer Client ID
	DefaultClientID = "clymene"
	// DefaultParallelism is the default parallelism for the metric processor
	DefaultParallelism = 1000
	// DefaultEncoding is the default metric encoding
	DefaultEncoding = kafka.EncodingProto
	// DefaultDeadlockInterval is the default deadlock interval
	DefaultDeadlockInterval = time.Duration(0)
)

// Options stores the configuration options for the ingester
type Options struct {
	kafkaConsumer.Configuration `mapstructure:",squash"`
	Parallelism                 int           `mapstructure:"parallelism"`
	Encoding                    string        `mapstructure:"encoding"`
	DeadlockInterval            time.Duration `mapstructure:"deadlock_interval"`
}

// AddFlags adds flags for Builder
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixBrokers,
		DefaultBroker,
		"The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234'")
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixPromtailTopic,
		DefaultPromtailTopic,
		"The name of the promtail kafka topic to consume from")
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixGroupID,
		DefaultGroupID,
		"The Consumer Group that clymene-ingester will be consuming on behalf of")
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixClientID,
		DefaultClientID,
		"The Consumer Client ID that clymene-ingester will use")
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixProtocolVersion,
		"",
		"Kafka protocol version - must be supported by kafka server")
	flagSet.String(
		KafkaConsumerConfigPrefix+SuffixEncoding,
		DefaultEncoding,
		fmt.Sprintf(`The encoding of metrics ("%s") consumed from kafka`, strings.Join(kafka.AllEncodings, "\", \"")))
	flagSet.String(
		ConfigPrefix+SuffixParallelism,
		strconv.Itoa(DefaultParallelism),
		"The number of messages to process in parallel")
	flagSet.Duration(
		ConfigPrefix+SuffixDeadlockInterval,
		DefaultDeadlockInterval,
		"Interval to check for deadlocks. If no messages gets processed in given time, clymene-ingester app will exit. Value of 0 disables deadlock check.")
	// Authentication flags
	auth.AddFlags(KafkaConsumerConfigPrefix, flagSet)
}

// InitFromViper initializes Builder with properties from viper
func (o *Options) InitFromViper(v *viper.Viper) {
	o.Brokers = strings.Split(stripWhiteSpace(v.GetString(KafkaConsumerConfigPrefix+SuffixBrokers)), ",")
	o.Topic = v.GetString(KafkaConsumerConfigPrefix + SuffixPromtailTopic)
	o.GroupID = v.GetString(KafkaConsumerConfigPrefix + SuffixGroupID)
	o.ClientID = v.GetString(KafkaConsumerConfigPrefix + SuffixClientID)
	o.ProtocolVersion = v.GetString(KafkaConsumerConfigPrefix + SuffixProtocolVersion)
	o.Encoding = v.GetString(KafkaConsumerConfigPrefix + SuffixEncoding)

	o.Parallelism = v.GetInt(ConfigPrefix + SuffixParallelism)
	o.DeadlockInterval = v.GetDuration(ConfigPrefix + SuffixDeadlockInterval)
	authenticationOptions := auth.AuthenticationConfig{}
	authenticationOptions.InitFromViper(KafkaConsumerConfigPrefix, v)
	o.AuthenticationConfig = authenticationOptions
}

// stripWhiteSpace removes all whitespace characters from a string
func stripWhiteSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}
