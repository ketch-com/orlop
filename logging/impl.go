package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/log"
	"time"
)

func New() (Logger, error) {
	return &loggerImpl{
		entry: log.New(),
	}, nil
}

type loggerImpl struct {
	entry *logrus.Entry
}

func (l loggerImpl) Printf(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l loggerImpl) Trace(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l loggerImpl) Tracef(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l loggerImpl) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l loggerImpl) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l loggerImpl) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l loggerImpl) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l loggerImpl) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l loggerImpl) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l loggerImpl) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l loggerImpl) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l loggerImpl) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l loggerImpl) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l loggerImpl) WithTime(t time.Time) Logger {
	return &loggerImpl{
		entry: l.entry.WithTime(t),
	}
}

func (l loggerImpl) WithContext(ctx context.Context) Logger {
	return &loggerImpl{
		entry: l.entry.WithContext(ctx),
	}
}

func (l loggerImpl) WithError(err error) Logger {
	return &loggerImpl{l.entry.WithError(err)}
}

func (l loggerImpl) WithField(key string, value interface{}) Logger {
	return &loggerImpl{l.entry.WithField(key, value)}
}

func (l loggerImpl) WithFields(fields ...interface{}) Logger {
	f := logrus.Fields{}
	for n := 0; n < len(fields) - 1; n += 2 {
		f[fields[n].(string)] = fields[n+1]
	}
	return &loggerImpl{l.entry.WithFields(f)}
}
