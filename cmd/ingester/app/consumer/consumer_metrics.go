package consumer

import (
	"strconv"

	"github.com/uber/jaeger-lib/metrics"
)

const consumerNamespace = "sarama-consumer"

type msgMetrics struct {
	counter     metrics.Counter
	offsetGauge metrics.Gauge
	lagGauge    metrics.Gauge
}

type errMetrics struct {
	errCounter metrics.Counter
}

type partitionMetrics struct {
	startCounter metrics.Counter
	closeCounter metrics.Counter
}

func (c *Consumer) namespace(partition int32) metrics.Factory {
	return c.metricsFactory.Namespace(metrics.NSOptions{Name: consumerNamespace, Tags: map[string]string{"partition": strconv.Itoa(int(partition))}})
}

func (c *Consumer) newMsgMetrics(partition int32) msgMetrics {
	f := c.namespace(partition)
	return msgMetrics{
		counter:     f.Counter(metrics.Options{Name: "messages", Tags: nil}),
		offsetGauge: f.Gauge(metrics.Options{Name: "current-offset", Tags: nil}),
		lagGauge:    f.Gauge(metrics.Options{Name: "offset-lag", Tags: nil}),
	}
}

func (c *Consumer) newErrMetrics(partition int32) errMetrics {
	return errMetrics{errCounter: c.namespace(partition).Counter(metrics.Options{Name: "errors", Tags: nil})}
}

func (c *Consumer) partitionMetrics(partition int32) partitionMetrics {
	f := c.namespace(partition)
	return partitionMetrics{
		closeCounter: f.Counter(metrics.Options{Name: "partition-close", Tags: nil}),
		startCounter: f.Counter(metrics.Options{Name: "partition-start", Tags: nil})}
}

func partitionsHeldGauge(metricsFactory metrics.Factory) metrics.Gauge {
	return metricsFactory.Namespace(metrics.NSOptions{Name: consumerNamespace, Tags: nil}).Gauge(metrics.Options{Name: "partitions-held", Tags: nil})
}
