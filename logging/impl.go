// Copyright (c) 2021 Ketch Kloud, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
