package log

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func FxLogger(entry *logrus.Entry) fx.Option {
	return fx.WithLogger(&fxLogger{
		entry: entry,
	})
}

type fxLogger struct {
	entry *logrus.Entry
}

func (l *fxLogger) LogEvent(event fxevent.Event) {
	if e, ok := event.(*fxevent.OnStartExecuting); ok {
		l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}).Debug("OnStartExecuting")
	} else if e, ok := event.(*fxevent.OnStartExecuted); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "method": e.Method, "runtime": e.Runtime}).Debug("OnStartExecuted")
	} else if e, ok := event.(*fxevent.OnStopExecuting); ok {
		l.entry.WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName}).Debug("OnStopExecuting")
	} else if e, ok := event.(*fxevent.OnStopExecuted); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"fn": e.FunctionName, "caller": e.CallerName, "runtime": e.Runtime}).Debug("OnStopExecuted")
	} else if e, ok := event.(*fxevent.Supplied); ok {
		l.addError(e.Err).WithField("caller", e.TypeName).Trace("Supplied")
	} else if e, ok := event.(*fxevent.Provided); ok {
		l.addError(e.Err).WithFields(logrus.Fields{"caller": e.ConstructorName, "outputTypeNames": e.OutputTypeNames}).Trace("Provided")
	} else if e, ok := event.(*fxevent.Invoking); ok {
		l.entry.WithField("fn", e.FunctionName).Debug("Invoking")
	} else if e, ok := event.(*fxevent.Invoked); ok {
		l.addError(e.Err).WithField("caller", e.FunctionName).Debug("Invoked")
	} else if e, ok := event.(*fxevent.Stopping); ok {
		l.entry.WithField("signal", e.Signal).Debug("Stopping")
	} else if e, ok := event.(*fxevent.Stopped); ok {
		l.addError(e.Err).Debug("Stopped")
	} else if e, ok := event.(*fxevent.RollingBack); ok {
		l.addError(e.StartErr).Debug("RollingBack")
	} else if e, ok := event.(*fxevent.RolledBack); ok {
		l.addError(e.Err).Debug("RolledBack")
	} else if e, ok := event.(*fxevent.Started); ok {
		l.addError(e.Err).Debug("Started")
	} else if e, ok := event.(*fxevent.LoggerInitialized); ok {
		l.addError(e.Err).WithField("constructor", e.ConstructorName).Debug("LoggerInitialized")
	}
}

func (l *fxLogger) addError(err error) *logrus.Entry {
	if err != nil {
		return l.entry.WithError(err)
	}

	return l.entry
}

