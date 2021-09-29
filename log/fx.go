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
		l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}).Trace("start executing")
	} else if e, ok := event.(*fxevent.OnStartExecuted); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "method": e.Method, "runtime": e.Runtime}).Trace("start executed")
	} else if e, ok := event.(*fxevent.OnStopExecuting); ok {
		l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}).Trace("stop executing")
	} else if e, ok := event.(*fxevent.OnStopExecuted); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "runtime": e.Runtime}).Trace("stop executed")
	} else if e, ok := event.(*fxevent.Supplied); ok {
		l.addError(e.Err).WithField("caller", e.TypeName).Trace("supplied")
	} else if e, ok := event.(*fxevent.Provided); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"caller": e.ConstructorName, "outputTypeNames": e.OutputTypeNames}).Trace("provided")
	} else if e, ok := event.(*fxevent.Invoking); ok {
		l.entry.WithField("fn", e.FunctionName).Trace("invoking")
	} else if e, ok := event.(*fxevent.Invoked); ok {
		l.addError(e.Err).WithField("caller", e.FunctionName).Trace("invoked")
	} else if e, ok := event.(*fxevent.Stopping); ok {
		l.entry.WithField("signal", e.Signal).Trace("stopping")
	} else if e, ok := event.(*fxevent.Stopped); ok {
		l.addError(e.Err).Trace("stopped")
	} else if e, ok := event.(*fxevent.RollingBack); ok {
		l.addError(e.StartErr).Trace("rolling back")
	} else if e, ok := event.(*fxevent.RolledBack); ok {
		l.addError(e.Err).Trace("rolled back")
	} else if e, ok := event.(*fxevent.Started); ok {
		l.addError(e.Err).Trace("started")
	} else if e, ok := event.(*fxevent.LoggerInitialized); ok {
		l.addError(e.Err).WithField("constructor", e.ConstructorName).Trace("logger initialized")
	}
}

func (l *fxLogger) addError(err error) *logrus.Entry {
	if err != nil {
		return l.entry.WithError(err)
	}

	return l.entry
}

