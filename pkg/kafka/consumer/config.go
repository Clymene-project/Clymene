package consumer

import (
	"github.com/Clymene-project/Clymene/pkg/kafka/auth"
	"go.uber.org/zap"
	"io"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// Consumer is an interface to features of Sarama that are necessary for the consumer
type Consumer interface {
	Partitions() <-chan cluster.PartitionConsumer
	MarkPartitionOffset(topic string, partition int32, offset int64, metadata string)
	io.Closer
}

// Builder builds a new kafka consumer
type Builder interface {
	NewConsumer() (Consumer, error)
}

// Configuration describes the configuration properties needed to create a Kafka consumer
type Configuration struct {
	auth.AuthenticationConfig `mapstructure:"authentication"`
	Consumer

	Brokers         []string `mapstructure:"brokers"`
	Topic           string   `mapstructure:"topic"`
	GroupID         string   `mapstructure:"group_id"`
	ClientID        string   `mapstructure:"client_id"`
	ProtocolVersion string   `mapstructure:"protocol_version"`
}

// NewConsumer creates a new kafka consumer
func (c *Configuration) NewConsumer(logger *zap.Logger) (Consumer, error) {
	saramaConfig := cluster.NewConfig()
	saramaConfig.Group.Mode = cluster.ConsumerModePartitions
	saramaConfig.ClientID = c.ClientID
	if len(c.ProtocolVersion) > 0 {
		ver, err := sarama.ParseKafkaVersion(c.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		saramaConfig.Config.Version = ver
	}
	if err := c.AuthenticationConfig.SetConfiguration(&saramaConfig.Config, logger); err != nil {
		return nil, err
	}
	return cluster.NewConsumer(c.Brokers, c.GroupID, []string{c.Topic}, saramaConfig)
}
