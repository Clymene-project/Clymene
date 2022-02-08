package util

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type LogAdapter struct {
	*zap.Logger
}

func NewLogAdapter(l *zap.Logger) LogAdapter {
	return LogAdapter{
		Logger: l,
	}
}

// Fatal implements tail.logger
func (l LogAdapter) Fatal(v ...interface{}) {
	l.Error(fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf implements tail.logger
func (l LogAdapter) Fatalf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
	os.Exit(1)
}

// Fatalln implements tail.logger
func (l LogAdapter) Fatalln(v ...interface{}) {
	l.Error(fmt.Sprint(v...))
	os.Exit(1)
}

// Panic implements tail.logger
func (l LogAdapter) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Error(s)
	panic(s)
}

// Panicf implements tail.logger
func (l LogAdapter) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...)
	l.Error(s)
	panic(s)
}

// Panicln implements tail.logger
func (l LogAdapter) Panicln(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Error(s)
	panic(s)
}

// Print implements tail.logger
func (l LogAdapter) Print(v ...interface{}) {
	l.Info(fmt.Sprint(v...))
}

// Printf implements tail.logger
func (l LogAdapter) Printf(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
}

// Println implements tail.logger
func (l LogAdapter) Println(v ...interface{}) {
	l.Info(fmt.Sprint(v...))
}

// TODO(dannyk): remove once weaveworks/common updates to go-kit/log
// 					-> we can then revert to using Level.Gokit
func LogFilter(l string) zapcore.Level {
	switch l {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
