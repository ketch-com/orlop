// Copyright (c) 2020 Ketch, Inc.
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

package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var loggerValue = struct{}{}

// ToContext adds a logger to the specified context, returning the new context
func ToContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerValue, logger)
}

// FromContext retrieves the logger for the provided context
func FromContext(ctx context.Context) *logrus.Entry {
	l := ctx.Value(loggerValue)
	if l != nil {
		if logger, ok := l.(*logrus.Entry); ok {
			return logger
		}
	}

	return WithContext(ctx)
}

// Writer returns a new PipeWriter
func Writer() *io.PipeWriter {
	return WriterLevel(logrus.GetLevel())
}

// WriterLevel returns a new PipeWriter
func WriterLevel(level logrus.Level) *io.PipeWriter {
	return logrus.StandardLogger().WriterLevel(level)
}

// New returns a new log Entry
func New() *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger())
}

// WithField returns a log Entry with the given Field
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// WithFields returns a log Entry with the given Fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

// WithError returns a log Entry with the given Error
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// WithContext returns a log Entry with the given Context
func WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

// WithTime returns a log Entry with the given Time
func WithTime(t time.Time) *logrus.Entry {
	return logrus.WithTime(t)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Debugln logs a debug message with a newline
func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Errorln logs a error message with a newline
func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

// Fatal logs a fatal message
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf logs a formatted fatal message
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Fatalln logs a fatal message with a newline
func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}

// Info logs an informational message
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof logs a formatted informational message
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Infoln logs an informational message with a newline
func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

// Trace logs a trace message
func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

// Tracef logs a formatted trace message
func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

// Traceln logs a trace message with a newline
func Traceln(args ...interface{}) {
	logrus.Traceln(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Warnln logs a warn message with a newline
func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
