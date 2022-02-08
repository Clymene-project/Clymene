package util

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// Provide an "event" interface for observability

// Temporary hack implementation to go via logger to stderr

var (
	// interface{} vars to avoid allocation on every call
	key   interface{} = "level" // masquerade as a level like debug, warn
	event interface{} = "event"

	eventLogger = zap.NewNop()
)

// Event is the log-like API for event sampling
func Event() *zap.Logger {
	return eventLogger
}

// InitEvents initializes event sampling, with the given frequency. Zero=off.
func InitEvents(freq int) {
	if freq <= 0 {
		eventLogger = zap.NewNop()
	} else {
		eventLogger = newEventLogger(freq)
	}
}

func newEventLogger(freq int) *zap.Logger {
	//l := zapcore.NewCore(os.Stderr)
	//
	//l = log.WithPrefix(l, key, event)
	//l
	//l = log.With(l, "ts", log.DefaultTimestampUTC)
	return zap.NewNop()
}

type samplingFilter struct {
	next  *zap.Logger
	freq  int
	count atomic.Int64
}

//func (e *samplingFilter) Log(keyvals ...interface{}) error {
//	count := e.count.Inc()
//	if count%int64(e.freㄷq) == 0 {
//		return e.next.Log(keyvals...)
//	}
//	return nil
//}
