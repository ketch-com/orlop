package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// Writer returns a new PipeWriter
func Writer() *io.PipeWriter {
	return logrus.New().WriterLevel(logrus.GetLevel())
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
	logrus.SetFormatter(&logrus.TextFormatter{})
}
