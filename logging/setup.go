package logging

import (
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/v2/env"
	stdlog "log"
)

// SetupLogging sets up logging for the environment and the default log level
func SetupLogging(env env.Environment, loglevel Level) {
	switch loglevel {
	case FatalLevel:
		logrus.SetLevel(logrus.FatalLevel)

	case ErrorLevel:
		logrus.SetLevel(logrus.ErrorLevel)

	case WarnLevel:
		logrus.SetLevel(logrus.WarnLevel)

	case InfoLevel:
		logrus.SetLevel(logrus.InfoLevel)

	case DebugLevel:
		logrus.SetLevel(logrus.DebugLevel)

	case TraceLevel:
		logrus.SetLevel(logrus.TraceLevel)

	default:
		if env.IsProduction() {
			logrus.SetLevel(logrus.WarnLevel)
		} else {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

	if env.IsLocal() {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		})
	}

	stdlog.SetOutput(logrus.New().Writer())
}
