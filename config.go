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
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

// Setup sets up the configuration system.
func Setup(serviceName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		viper.AutomaticEnv()
		viper.SetEnvPrefix(strings.ToUpper(serviceName))
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

		env := Env()

		envFiles := []string{".env"}
		if env.IsLocal() {
			envFiles = append(envFiles, ".env.local")
		} else {
			envFiles = append(envFiles, ".env."+env.String())
		}

		_ = godotenv.Overload(envFiles...)

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			return err
		}

		configFile := cmd.Flags().Lookup("config").Value.String()
		viper.SetConfigFile(configFile)
		viper.ReadInConfig()

		switch viper.GetString("loglevel") {
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

		return nil
	}
}

// Lifted from viper
func stringToStringConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// An empty string would cause an empty map
	if len(val) == 0 {
		return map[string]string{}, nil
	}
	r := csv.NewReader(strings.NewReader(val))
	ss, err := r.Read()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		}
		out[kv[0]] = kv[1]
	}
	return out, nil
}

func decodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.String() == "string" && t.String() == "map[string]string" {
		return stringToStringConv(data.(string))
	}

	return data, nil
}

// Unmarshal returns configuration in the specified object.
func Unmarshal(v interface{}) error {
	return viper.Unmarshal(v, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		decodeHook,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
}

// MakeCommandKeyPrefix returns a function that prepends the given prefix to the key
func MakeCommandKeyPrefix(prefix []string) func(key string) string {
	return func(key string) string {
		if len(prefix) == 0 {
			return key
		}

		return strings.Join(prefix, ".") + "." + key
	}
}
