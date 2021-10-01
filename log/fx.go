// Copyright (c) 2020 Ketch Kloud, Inc.
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
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func FxLogger(entry *logrus.Entry) fx.Option {
	return fx.WithLogger(func () fxevent.Logger {
		return &fxLogger{
			entry: entry,
		}
	})
}

type fxLogger struct {
	entry *logrus.Entry
}

func (l *fxLogger) LogEvent(event fxevent.Event) {
	if e, ok := event.(*fxevent.OnStartExecuting); ok {
		l.message(nil, "start executing", l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}))
	} else if e, ok := event.(*fxevent.OnStartExecuted); ok {
		l.message(e.Err, "start executed", l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "method": e.Method, "runtime": e.Runtime}))
	} else if e, ok := event.(*fxevent.OnStopExecuting); ok {
		l.message(nil, "stop executing", l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}))
	} else if e, ok := event.(*fxevent.OnStopExecuted); ok {
		l.message(e.Err, "stop executed", l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "runtime": e.Runtime}))
	} else if e, ok := event.(*fxevent.Supplied); ok {
		l.message(e.Err, "supplied", l.entry.WithField("caller", e.TypeName))
	} else if e, ok := event.(*fxevent.Provided); ok {
		l.message(e.Err, "provided", l.entry.WithFields(logrus.Fields{"caller": e.ConstructorName, "outputTypeNames": e.OutputTypeNames}))
	} else if e, ok := event.(*fxevent.Invoking); ok {
		l.message(nil, "invoking", l.entry.WithField("fn", e.FunctionName))
	} else if e, ok := event.(*fxevent.Invoked); ok {
		l.message(e.Err, "invoked", l.entry.WithField("caller", e.FunctionName))
	} else if e, ok := event.(*fxevent.Stopping); ok {
		l.message(nil, "stopping", l.entry.WithField("signal", e.Signal))
	} else if e, ok := event.(*fxevent.Stopped); ok {
		l.message(e.Err, "stopped", l.entry)
	} else if e, ok := event.(*fxevent.RollingBack); ok {
		l.message(e.StartErr, "rolling back", l.entry)
	} else if e, ok := event.(*fxevent.RolledBack); ok {
		l.message(e.Err, "rolled back", l.entry)
	} else if e, ok := event.(*fxevent.Started); ok {
		l.message(e.Err, "started", l.entry)
	} else if e, ok := event.(*fxevent.LoggerInitialized); ok {
		l.message(e.Err, "logger initialized", l.entry.WithField("constructor", e.ConstructorName))
	}
}

func (l *fxLogger) message(err error, message string, log *logrus.Entry) {
	if err != nil {
		log.WithError(err).Errorln(message)
	}

	log.Trace(message)
}
