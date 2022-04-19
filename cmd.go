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

package orlop

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/log"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
	stdlog "log"
	"os"
	"reflect"
	"sort"
	"strings"
)

// Runner represents a command runner
type Runner struct {
	prefix      string
	prevPreRunE func(cmd *cobra.Command, args []string) error
}

// NewRunner creates a new Runner
func NewRunner(prefix string) *Runner {
	return &Runner{
		prefix: prefix,
	}
}

// SetupRoot sets up the root Command
func (r *Runner) SetupRoot(cmd *cobra.Command) *Runner {
	if cmd.PersistentFlags().Lookup("env") == nil {
		cmd.PersistentFlags().String("env", strings.ToLower(r.Getenv("environment")), "specifies the environment")
	}
	if cmd.PersistentFlags().Lookup("loglevel") == nil {
		cmd.PersistentFlags().String("loglevel", strings.ToLower(r.Getenv("loglevel")), "specifies the log level")
	}
	if cmd.PersistentFlags().Lookup("config") == nil {
		cmd.PersistentFlags().StringSlice("config", nil, "specifies a .env configuration file to load")
	}

	r.prevPreRunE = cmd.PersistentPreRunE

	if r.prevPreRunE == nil {
		r.prevPreRunE = func(cmd *cobra.Command, args []string) error {
			return nil
		}
	}

	cmd.PersistentPreRunE = r.preRunE
	return r
}

// Setup sets up the Command
func (r *Runner) Setup(cmd *cobra.Command, runner interface{}) *Runner {
	if cmd.RunE == nil {
		cmd.RunE = r.runE(runner)
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "output the config environment variables and exits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config Manager List ALL

			vars, err := GetVariablesFromConfig(r.prefix, &struct{}{})
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

			return nil
		},
	})

	return r
}

func (r *Runner) preRunE(cmd *cobra.Command, args []string) error {
	envFlag, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	loglevelFlag, err := cmd.Flags().GetString("loglevel")
	if err != nil {
		return err
	}

	configFiles, err := cmd.Flags().GetStringSlice("config")
	if err != nil {
		return err
	}

	// First figure out the environment
	env := Environment(envFlag)

	// Load the environment from files
	LoadEnvironment(env, configFiles...)

	// Setup logging
	r.SetupLogging(env, loglevelFlag)

	return r.prevPreRunE(cmd, args)
}

func (r *Runner) runE(runner interface{}) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		envFlag, err := cmd.Flags().GetString("env")
		if err != nil {
			return err
		}

		ctx, span := tracer.Start(cmd.Context(), "Run")
		defer span.End()

		// First figure out the environment
		span.SetAttributes(attribute.String("env", Environment(envFlag).String()))
		span.SetAttributes(semconv.ServiceNameKey.String(r.prefix))

		l := log.New()

		loglevelFlag, err := cmd.Flags().GetString("loglevel")
		if err != nil {
			return err
		}

		if module, ok := runner.(fx.Option); ok {
			runner = func(ctx context.Context) error {
				app := fx.New(
					logging.WithLogger(l),
					FxContext(ctx),
					fx.Supply(cmd),
					fx.Supply(service.Name(r.prefix)),
					fx.Supply(logging.Level(loglevelFlag)),
					Module,
					module,
				)

				//populate() error
				//Load() error

				app.Run()

				return app.Err()
			}
		}

		// Call the runner
		out := reflect.ValueOf(runner).Call([]reflect.Value{
			reflect.ValueOf(log.ToContext(ctx, l)),
		})

		// Handle any result
		if len(out) > 0 && out[0].IsValid() && !out[0].IsNil() {
			e := out[0].MethodByName("Error")
			out = e.Call([]reflect.Value{})
			if len(out) > 0 && out[0].IsValid() {
				return errors.New(out[0].String())
			}
		}

		return nil
	}
}

// Getenv returns the value of the environment variabled named `key`
func (r *Runner) Getenv(key string) string {
	return os.Getenv(strcase.ToScreamingSnake(strings.Join([]string{r.prefix, key}, "_")))
}

// SetupLogging sets up logging for the environment and the default log level
func (r *Runner) SetupLogging(env Environment, loglevel string) {
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
}

// Run loads config and then executes the given runner
func Run(prefix string, runner interface{}) {
	var cmd = &cobra.Command{
		Use:              prefix,
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	NewRunner(prefix).SetupRoot(cmd).Setup(cmd, runner)

	if err := cmd.Execute(); err != nil {
		log.WithError(err).Fatal(err)
	}
}
