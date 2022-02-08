package kafka

import (
	"context"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/grafana/dskit/backoff"
)

var defaultBackOff = backoff.Config{
	MinBackoff: 1 * time.Second,
	MaxBackoff: 60 * time.Second,
	MaxRetries: 20,
}

type RunnableTarget interface {
	target.Target
	run()
}

type TargetDiscoverer interface {
	NewTarget(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) (RunnableTarget, error)
}

// consumer handle a group consumer instance.
// It will create a new target for every consumer claim using the `TargetDiscoverer`.
type consumer struct {
	sarama.ConsumerGroup
	discoverer TargetDiscoverer
	logger     *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mutex          sync.Mutex // used during rebalancing setup and tear down
	activeTargets  []target.Target
	droppedTargets []target.Target
}

// start starts the consumer for a given list of topics.
func (c *consumer) start(ctx context.Context, topics []string) {
	c.wg.Wait()
	c.wg.Add(1)

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.logger.Info("starting consumer", zap.String("topics", fmt.Sprintf("%+v", topics)))

	go func() {
		defer c.wg.Done()
		backoff := backoff.New(c.ctx, defaultBackOff)
		for {
			// Calling Consume in an infinite loop in case rebalancing is kicking in.
			// In which case all claims will be renewed.
			err := c.ConsumerGroup.Consume(c.ctx, topics, c)
			if err != nil && err != context.Canceled {
				c.logger.Error("error from the consumer, retrying...", zap.Error(err))
				// backoff before re-trying.
				backoff.Wait()
				if backoff.Ongoing() {
					continue
				}
				c.logger.Error("maximun error from the consumer reached", zap.String("last_err", err.Error()))
				return
			}
			if c.ctx.Err() != nil || err == context.Canceled {
				c.logger.Info("stopping consumer", zap.String("topics", fmt.Sprintf("%+v", topics)))
				return
			}
			backoff.Reset()
		}
	}()
}

// ConsumeClaim creates a target for the given received claim and start reading message from it.
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.wg.Add(1)
	defer c.wg.Done()

	t, err := c.discoverer.NewTarget(session, claim)
	if err != nil {
		return err
	}
	if len(t.Labels()) == 0 {
		c.addDroppedTarget(t)
		t.run()
		return nil
	}
	c.addTarget(t)
	c.logger.Info("consuming topic", zap.Any("details", t.Details()))
	t.run()

	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(session sarama.ConsumerGroupSession) error {
	c.resetTargets()
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.resetTargets()
	return nil
}

// stop stops the consumer.
func (c *consumer) stop() {
	c.cancel()
	c.wg.Wait()
	c.resetTargets()
}

func (c *consumer) resetTargets() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.activeTargets = nil
	c.droppedTargets = nil
}

func (c *consumer) getActiveTargets() []target.Target {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.activeTargets
}

func (c *consumer) getDroppedTargets() []target.Target {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.droppedTargets
}

func (c *consumer) addTarget(t target.Target) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.activeTargets = append(c.activeTargets, t)
}

func (c *consumer) addDroppedTarget(t target.Target) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.droppedTargets = append(c.droppedTargets, t)
}
