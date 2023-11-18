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

package cmd

import (
	"context"
	"fmt"
	stdlog "log"
	"sort"
	"strings"

	"go.ketch.com/lib/orlop/v2"
	"go.ketch.com/lib/orlop/v2/env"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/log"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
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
func (r *Runner) Setup(cmd *cobra.Command, options ...fx.Option) *Runner {
	if cmd.RunE == nil {
		cmd.RunE = r.runE(options...)
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "output the config environment variables and exits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfgMgr config.Provider
			app := fx.New(
				fx.NopLogger,
				fx.Provide(func() context.Context { return cmd.Context() }),
				fx.Supply(cmd),
				fx.Supply(service.Name(r.prefix)),
				fx.Supply(logging.FatalLevel),
				orlop.Module,
				fx.Options(options...),
				fx.Populate(&cfgMgr),
			)

			if err := app.Start(context.Background()); err != nil {
				panic(err)
			}
			defer app.Stop(context.Background())

			vars, err := cfgMgr.List(cmd.Context())
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
	e := env.Environment(envFlag)

	// Load the environment from files
	e.Load(configFiles...)

	// Setup logging
	r.SetupLogging(e, loglevelFlag)

	return r.prevPreRunE(cmd, args)
}

func (r *Runner) runE(options ...fx.Option) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		l := log.New()

		loglevelFlag, err := cmd.Flags().GetString("loglevel")
		if err != nil {
			return err
		}

		app := fx.New(
			logging.WithLogger(l),
			fx.Provide(func() context.Context { return cmd.Context() }),
			fx.Supply(cmd),
			fx.Supply(service.Name(r.prefix)),
			fx.Supply(logging.Level(loglevelFlag)),
			orlop.Module,
			fx.Options(options...),
		)

		app.Run()

		return app.Err()
	}
}

// Getenv returns the value of the environment variable named `key`
func (r *Runner) Getenv(key string) string {
	return config.GetEnv(r.prefix, key)
}

// SetupLogging sets up logging for the environment and the default log level
func (r *Runner) SetupLogging(e env.Environment, loglevel string) {
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
		if e.IsProduction() {
			logrus.SetLevel(logrus.WarnLevel)
		} else {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

	if e.IsLocal() {
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
func Run(prefix string, module fx.Option) {
	var cmd = &cobra.Command{
		Use:              prefix,
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	NewRunner(prefix).SetupRoot(cmd).Setup(cmd, module)

	if err := cmd.Execute(); err != nil {
		log.WithError(err).Fatal(err)
	}
}
