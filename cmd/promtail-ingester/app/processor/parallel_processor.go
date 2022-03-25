package processor

import (
	"sync"

	"go.uber.org/zap"
)

// ParallelProcessor is a processor that processes in parallel using a pool of goroutines
type ParallelProcessor struct {
	messages    chan Message
	processor   LogProcessor
	numRoutines int

	logger *zap.Logger
	closed chan struct{}
	wg     sync.WaitGroup
}

// NewParallelProcessor creates a new parallel processor
func NewParallelProcessor(
	processor LogProcessor,
	parallelism int,
	logger *zap.Logger) *ParallelProcessor {
	return &ParallelProcessor{
		logger:      logger,
		messages:    make(chan Message),
		processor:   processor,
		numRoutines: parallelism,
		closed:      make(chan struct{}),
	}
}

// Start begins processing queued messages
func (k *ParallelProcessor) Start() {
	k.logger.Debug("Spawning goroutines to process messages", zap.Int("num_routines", k.numRoutines))
	for i := 0; i < k.numRoutines; i++ {
		k.wg.Add(1)
		go func() {
			for {
				select {
				case msg := <-k.messages:
					k.processor.Process(msg)
				case <-k.closed:
					k.wg.Done()
					return
				}
			}
		}()
	}
}

// Process queues a message for processing
func (k *ParallelProcessor) Process(message Message) error {
	k.messages <- message
	return nil
}

// Close terminates all running goroutines
func (k *ParallelProcessor) Close() error {
	k.logger.Debug("Initiated shutdown of processor goroutines")
	close(k.closed)
	k.wg.Wait()
	k.logger.Info("Completed shutdown of processor goroutines")
	return nil
}
