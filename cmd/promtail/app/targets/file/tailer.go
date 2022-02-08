package file

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/pkg/logproto"
	util "github.com/Clymene-project/Clymene/pkg/lokiutil"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/prometheus/common/model"
	"go.uber.org/atomic"
)

type tailer struct {
	metrics   *Metrics
	logger    *zap.Logger
	handler   api.EntryHandler
	positions positions.Positions

	path string
	tail *tail.Tail

	posAndSizeMtx sync.Mutex
	stopOnce      sync.Once

	running *atomic.Bool
	posquit chan struct{}
	posdone chan struct{}
	done    chan struct{}
}

func newTailer(metrics *Metrics, logger *zap.Logger, handler api.EntryHandler, positions positions.Positions, path string) (*tailer, error) {
	// Simple check to make sure the file we are tailing doesn't
	// have a position already saved which is past the end of the file.
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	pos, err := positions.Get(path)
	if err != nil {
		return nil, err
	}

	if fi.Size() < pos {
		positions.Remove(path)
	}

	tail, err := tail.TailFile(path, tail.Config{
		Follow:    true,
		Poll:      true,
		ReOpen:    true,
		MustExist: true,
		Location: &tail.SeekInfo{
			Offset: pos,
			Whence: 0,
		},
		Logger: util.NewLogAdapter(logger),
	})
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("component", "tailer"))
	tailer := &tailer{
		metrics:   metrics,
		logger:    logger,
		handler:   api.AddLabelsMiddleware(model.LabelSet{FilenameLabel: model.LabelValue(path)}).Wrap(handler),
		positions: positions,
		path:      path,
		tail:      tail,
		running:   atomic.NewBool(false),
		posquit:   make(chan struct{}),
		posdone:   make(chan struct{}),
		done:      make(chan struct{}),
	}

	go tailer.readLines()
	go tailer.updatePosition()
	metrics.filesActive.Add(1.)
	return tailer, nil
}

// updatePosition is run in a goroutine and checks the current size of the file and saves it to the positions file
// at a regular interval. If there is ever an error it stops the tailer and exits, the tailer will be re-opened
// by the filetarget sync method if it still exists and will start reading from the last successful entry in the
// positions file.
func (t *tailer) updatePosition() {
	positionSyncPeriod := t.positions.SyncPeriod()
	positionWait := time.NewTicker(positionSyncPeriod)
	defer func() {
		positionWait.Stop()
		t.logger.Info("position timer: exited", zap.String("path", t.path))
		close(t.posdone)
	}()

	for {
		select {
		case <-positionWait.C:
			err := t.markPositionAndSize()
			if err != nil {
				t.logger.Error("position timer: error getting tail position and/or size, stopping tailer", zap.String("path", t.path), zap.Error(err))
				err := t.tail.Stop()
				if err != nil {
					t.logger.Error("position timer: error stopping tailer", zap.String("path", t.path), zap.Error(err))
				}
				return
			}
		case <-t.posquit:
			return
		}
	}
}

// readLines runs in a goroutine and consumes the t.tail.Lines channel from the underlying tailer.
// it will only exit when that channel is closed. This is important to avoid a deadlock in the underlying
// tailer which can happen if there are unread lines in this channel and the Stop method on the tailer
// is called, the underlying tailer will never exit if there are unread lines in the t.tail.Lines channel
func (t *tailer) readLines() {
	t.logger.Info("tail routine: started", zap.String("path", t.path))

	t.running.Store(true)

	// This function runs in a goroutine, if it exits this tailer will never do any more tailing.
	// Clean everything up.
	defer func() {
		t.cleanupMetrics()
		t.running.Store(false)
		t.logger.Info("tail routine: exited", zap.String("path", t.path))
		close(t.done)
	}()
	entries := t.handler.Chan()
	for {
		line, ok := <-t.tail.Lines
		if !ok {
			t.logger.Info("tail routine: tail channel closed, stopping tailer", zap.String("path", t.path), zap.String("reason", t.tail.Tomb.Err()))
			return
		}

		// Note currently the tail implementation hardcodes Err to nil, this should never hit.
		if line.Err != nil {
			t.logger.Error("tail routine: error reading line", zap.String("path", t.path), zap.String("error", line.Err))
			continue
		}

		t.metrics.readLines.WithLabelValues(t.path).Inc()
		t.metrics.logLengthHistogram.WithLabelValues(t.path).Observe(float64(len(line.Text)))
		entries <- api.Entry{
			Labels: model.LabelSet{},
			Entry: logproto.Entry{
				Timestamp: line.Time,
				Line:      line.Text,
			},
		}

	}
}

func (t *tailer) markPositionAndSize() error {
	// Lock this update as there are 2 timers calling this routine, the sync in filetarget and the positions sync in this file.
	t.posAndSizeMtx.Lock()
	defer t.posAndSizeMtx.Unlock()

	size, err := t.tail.Size()
	if err != nil {
		// If the file no longer exists, no need to save position information
		if err == os.ErrNotExist {
			t.logger.Info("skipping update of position for a file which does not currently exist", zap.String("path", t.path))
			return nil
		}
		return err
	}
	t.metrics.totalBytes.WithLabelValues(t.path).Set(float64(size))

	pos, err := t.tail.Tell()
	if err != nil {
		return err
	}
	t.metrics.readBytes.WithLabelValues(t.path).Set(float64(pos))
	t.positions.Put(t.path, pos)

	return nil
}

func (t *tailer) stop() {
	// stop can be called by two separate threads in filetarget, to avoid a panic closing channels more than once
	// we wrap the stop in a sync.Once.
	t.stopOnce.Do(func() {
		// Shut down the position marker thread
		close(t.posquit)
		<-t.posdone

		// Save the current position before shutting down tailer
		err := t.markPositionAndSize()
		if err != nil {
			t.logger.Error("error marking file position when stopping tailer", zap.String("path", t.path), zap.Error(err))
		}

		// Stop the underlying tailer
		err = t.tail.Stop()
		if err != nil {
			t.logger.Error("error stopping tailer", zap.String("path", t.path), zap.Error(err))
		}
		// Wait for readLines() to consume all the remaining messages and exit when the channel is closed
		<-t.done
		t.logger.Info("stopped tailing file", zap.String("path", t.path))
		t.handler.Stop()
	})
}

func (t *tailer) isRunning() bool {
	return t.running.Load()
}

// cleanupMetrics removes all metrics exported by this tailer
func (t *tailer) cleanupMetrics() {
	// When we stop tailing the file, also un-export metrics related to the file
	t.metrics.filesActive.Add(-1.)
	t.metrics.readLines.DeleteLabelValues(t.path)
	t.metrics.readBytes.DeleteLabelValues(t.path)
	t.metrics.totalBytes.DeleteLabelValues(t.path)
	t.metrics.logLengthHistogram.DeleteLabelValues(t.path)
}
