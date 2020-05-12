// Copyright (c) 2020 SwitchBit, Inc.
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
	"log"
	"os"
	"reflect"
	"strings"
)

func Run(prefix string, runner interface{}, cfg interface{}) {
	if len(os.Args) > 1 {
		if os.Args[1] == "--init" {
			vars, err := GetVariablesFromConfig(prefix, cfg)
			if err != nil {
				logrus.Fatal(err)
			}

			fmt.Println(strings.Join(vars, "\n"))
			return
		}
	}

	// First figure out the environment
	env := Environment(strings.ToLower(getenv(prefix, "environment")))

	// Load the environment from files
	loadEnvironment(env)

	// Setup logging
	setupLogging(prefix, env)

	// Unmarshal the configuration
	err := Unmarshal(prefix, cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	// Call the runner
	out := reflect.ValueOf(runner).Call([]reflect.Value{
		reflect.ValueOf(context.TODO()),
		reflect.ValueOf(cfg),
	})

	// Handle any result
	if len(out) > 0 && out[0].IsValid() {
		e := out[0].MethodByName("Error")
		out = e.Call([]reflect.Value{})
	}

	if len(out) > 0 && out[0].IsValid() {
		logrus.Fatal(out[0].String())
	}
}

func getenv(prefix string, key string) string {
	return os.Getenv(strcase.ToScreamingSnake(strings.Join([]string{prefix, key}, "_")))
}

func setupLogging(prefix string, env Environment) {
	switch strings.ToLower(getenv(prefix, "loglevel")) {
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
			ForceColors: true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	log.SetOutput(logrus.New().Writer())
}
