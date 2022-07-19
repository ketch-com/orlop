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
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/log"
	"go.ketch.com/lib/orlop/v2/request"
	"time"
)

func New(env env.Environment, loglevel Level) (Logger, error) {
	SetupLogging(env, loglevel)

	return &loggerImpl{
		entry: log.New(),
	}, nil
}

type loggerImpl struct {
	entry *logrus.Entry
}

func (l loggerImpl) Printf(format string, args ...any) {
	l.entry.Infof(format, args...)
}

func (l loggerImpl) Trace(args ...any) {
	l.entry.Debug(args...)
}

func (l loggerImpl) Tracef(format string, args ...any) {
	l.entry.Debugf(format, args...)
}

func (l loggerImpl) Debug(args ...any) {
	l.entry.Debug(args...)
}

func (l loggerImpl) Debugf(format string, args ...any) {
	l.entry.Debugf(format, args...)
}

func (l loggerImpl) Info(args ...any) {
	l.entry.Info(args...)
}

func (l loggerImpl) Infof(format string, args ...any) {
	l.entry.Infof(format, args...)
}

func (l loggerImpl) Warn(args ...any) {
	l.entry.Warn(args...)
}

func (l loggerImpl) Warnf(format string, args ...any) {
	l.entry.Warnf(format, args...)
}

func (l loggerImpl) Error(args ...any) {
	l.entry.Error(args...)
}

func (l loggerImpl) Errorf(format string, args ...any) {
	l.entry.Errorf(format, args...)
}

func (l loggerImpl) Fatal(args ...any) {
	l.entry.Fatal(args...)
}

func (l loggerImpl) Fatalf(format string, args ...any) {
	l.entry.Fatalf(format, args...)
}

func (l loggerImpl) WithTime(t time.Time) Logger {
	return &loggerImpl{
		entry: l.entry.WithTime(t),
	}
}

func (l loggerImpl) WithContext(ctx context.Context) Logger {
	entry := l.entry.WithContext(ctx)

	for k, v := range request.Values(ctx) {
		entry = entry.WithField(k, v)
	}

	return &loggerImpl{
		entry: entry,
	}
}

func (l loggerImpl) WithError(err error) Logger {
	return &loggerImpl{l.entry.WithError(err)}
}

func (l loggerImpl) WithField(key string, value any) Logger {
	return &loggerImpl{l.entry.WithField(key, value)}
}

func (l loggerImpl) WithFields(fields ...any) Logger {
	f := logrus.Fields{}
	for n := 0; n < len(fields)-1; n += 2 {
		f[fields[n].(string)] = fields[n+1]
	}
	return &loggerImpl{l.entry.WithFields(f)}
}
