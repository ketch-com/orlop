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
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/env"
	"go.uber.org/fx"
	stdlog "log"
)

var Module = fx.Options(
	fx.Invoke(SetupLogging),
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
