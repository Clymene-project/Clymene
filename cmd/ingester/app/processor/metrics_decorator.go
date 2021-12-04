package processor

import (
	"io"
	"time"

	"github.com/uber/jaeger-lib/metrics"
)

type metricsDecorator struct {
	errors    metrics.Counter
	latency   metrics.Timer
	processor MetricProcessor
	io.Closer
}

// NewDecoratedProcessor returns a processor with metrics
func NewDecoratedProcessor(f metrics.Factory, processor MetricProcessor) MetricProcessor {
	m := f.Namespace(metrics.NSOptions{Name: "metric-processor", Tags: nil})
	return &metricsDecorator{
		errors:    m.Counter(metrics.Options{Name: "errors", Tags: nil}),
		latency:   m.Timer(metrics.TimerOptions{Name: "latency", Tags: nil}),
		processor: processor,
	}
}

func (d *metricsDecorator) Process(message Message) error {
	now := time.Now()

	err := d.processor.Process(message)
	d.latency.Record(time.Since(now))
	if err != nil {
		d.errors.Inc(1)
	}
	return err
}
