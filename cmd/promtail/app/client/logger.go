package client

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/fatih/color"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"os"
	"runtime"
	"sync"
	"text/tabwriter"
)

var (
	yellow = color.New(color.FgYellow)
	blue   = color.New(color.FgBlue)
)

func init() {
	if runtime.GOOS == "windows" {
		yellow.DisableColor()
		blue.DisableColor()
	}
}

type logger struct {
	*tabwriter.Writer
	sync.Mutex
	entries chan api.Entry

	once sync.Once
}

// NewLogger creates a new client logger that logs entries instead of sending them.
func NewLogger(options Options, logWriter logstore.Writer, factory metrics.Factory, log *zap.Logger) (Client, error) {
	// make sure the clients config is valid
	c, err := NewMulti(options, logWriter, factory, log)
	if err != nil {
		return nil, err
	}
	c.Stop()

	entries := make(chan api.Entry)
	l := &logger{
		Writer:  tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0),
		entries: entries,
	}
	go l.run()
	return l, nil
}

func (l *logger) Stop() {
	l.once.Do(func() { close(l.entries) })
}

func (l *logger) Chan() chan<- api.Entry {
	return l.entries
}

func (l *logger) run() {
	for e := range l.entries {
		fmt.Fprint(l.Writer, blue.Sprint(e.Timestamp.Format("2006-01-02T15:04:05.999999999-0700")))
		fmt.Fprint(l.Writer, "\t")
		fmt.Fprint(l.Writer, yellow.Sprint(e.Labels.String()))
		fmt.Fprint(l.Writer, "\t")
		fmt.Fprint(l.Writer, e.Line)
		fmt.Fprint(l.Writer, "\n")
		l.Flush()
	}
}
func (l *logger) StopNow() { l.Stop() }
