package client

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"sync"
)

// MultiClient is client pushing to one or more loki instances.
type MultiClient struct {
	client  Client
	entries chan api.Entry
	wg      sync.WaitGroup

	once sync.Once
}

// NewMulti creates a new client
func NewMulti(options Options, logWriter logstore.Writer, factory metrics.Factory, log *zap.Logger) (Client, error) {
	client, err := New(options, logWriter, factory, log)
	if err != nil {
		return nil, err
	}
	multi := &MultiClient{
		client:  client,
		entries: make(chan api.Entry),
	}
	multi.start()
	return multi, nil
}

func (m *MultiClient) start() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for e := range m.entries {
			m.client.Chan() <- e
		}
	}()
}

func (m *MultiClient) Chan() chan<- api.Entry {
	return m.entries
}

// Stop implements Client
func (m *MultiClient) Stop() {
	m.once.Do(func() { close(m.entries) })
	m.wg.Wait()
	m.client.Stop()
}

// StopNow implements Client
func (m *MultiClient) StopNow() {
	m.client.StopNow()
}
