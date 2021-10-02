package logging

import (
	"context"
	"time"
)

type Logger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Printf(format string, args ...interface{})

	WithTime(t time.Time) Logger
	WithContext(ctx context.Context) Logger
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields ...interface{}) Logger
}
