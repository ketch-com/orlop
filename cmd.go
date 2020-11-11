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

package orlop

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.ketch.com/lib/orlop/log"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/label"
	stdlog "log"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

// Run loads config and then executes the given runner
func Run(prefix string, runner interface{}, cfg interface{}) {
	var configFiles []string
	var initFlag bool
	var envFlag string
	var loglevelFlag string
	pflag.BoolVar(&initFlag, "init", false, "outputs the config environment variables and exits")
	pflag.StringVar(&envFlag, "env", strings.ToLower(getenv(prefix, "environment")), "specifies the environment")
	pflag.StringVar(&loglevelFlag, "loglevel", strings.ToLower(getenv(prefix, "loglevel")), "specifies the log level")
	pflag.StringSliceVar(&configFiles, "config", nil, "specifies a .env configuration file to load")
	pflag.Parse()

	if initFlag {
		vars, err := GetVariablesFromConfig(prefix, cfg)
		if err != nil {
			log.WithError(err).Fatal("could not create variables")
		}

		sort.Strings(vars)
		for _, v := range vars {
			if strings.Contains(v, "=#") {
				fmt.Println("#" + v)
			} else {
				fmt.Println(v)
			}
		}
		return
	}

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		log.WithError(err).Fatal("could not start runtime tracing")
	}

	ctx, span := tracer.Start(context.Background(), "Run")
	defer span.End()

	// First figure out the environment
	env := Environment(envFlag)
	span.SetAttributes(label.String("env", env.String()))

	// Load the environment from files
	loadEnvironment(env, configFiles...)

	// Setup logging
	ctx = setupLogging(ctx, env, loglevelFlag)

	// Unmarshal the configuration
	err := Unmarshal(prefix, cfg)
	if err != nil {
		log.FromContext(ctx).Fatal(err)
	}

	// Call the runner
	out := reflect.ValueOf(runner).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(cfg),
	})

	// Handle any result
	if len(out) > 0 && out[0].IsValid() && !out[0].IsNil() {
		e := out[0].MethodByName("Error")
		out = e.Call([]reflect.Value{})
		if len(out) > 0 && out[0].IsValid() {
			log.FromContext(ctx).Fatal(out[0].String())
		}
	}
}

func getenv(prefix string, key string) string {
	return os.Getenv(strcase.ToScreamingSnake(strings.Join([]string{prefix, key}, "_")))
}

func setupLogging(ctx context.Context, env Environment, loglevel string) context.Context {
	switch loglevel {
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)

	case "error":
		logrus.SetLevel(logrus.ErrorLevel)

	case "info":
		logrus.SetLevel(logrus.InfoLevel)

	case "debug":
		logrus.SetLevel(logrus.DebugLevel)

	case "trace":
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

	return log.ToContext(ctx, log.New())
}
