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

package config

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.ketch.com/lib/orlop/v2/errors"
)

type fieldSetter func(value reflect.Value, input string) error

var knownSetters map[string]fieldSetter

// RegisterConfigParser registers a config parser
func RegisterConfigParser(typeName string, parser fieldSetter) {
	knownSetters[typeName] = parser
}

func unmarshalTextSetter(value reflect.Value, input string) error {
	m := value.Addr().MethodByName("UnmarshalText")
	m.Call([]reflect.Value{reflect.ValueOf([]byte(input))})
	return nil
}

func unmarshalJSONSetter(value reflect.Value, input string) error {
	if len(input) > 0 {
		m := value.Addr().MethodByName("UnmarshalJSON")
		m.Call([]reflect.Value{reflect.ValueOf([]byte(input))})
	}
	return nil
}

func boolFieldSetter(value reflect.Value, input string) error {
	value.SetBool(strings.ToLower(input) == "true")
	return nil
}

func intFieldSetter(value reflect.Value, input string) error {
	if len(input) == 0 {
		input = "0"
	}

	i, err := strconv.ParseInt(input, 0, 0)
	if err != nil {
		return errors.Errorf("could not parse '%s' as integer", input)
	}

	value.SetInt(i)
	return nil
}

func uintFieldSetter(value reflect.Value, input string) error {
	if len(input) == 0 {
		input = "0"
	}

	i, err := strconv.ParseUint(input, 0, 0)
	if err != nil {
		return errors.Errorf("could not parse '%s' as integer", input)
	}

	value.SetUint(i)
	return nil
}

func floatFieldSetter(value reflect.Value, input string) error {
	if len(input) == 0 {
		input = "0"
	}

	i, err := strconv.ParseFloat(input, 0)
	if err != nil {
		return errors.Errorf("could not parse '%s' as float", value)
	}

	value.SetFloat(i)
	return nil
}

func mapFieldSetter(value reflect.Value, input string) error {
	m := reflect.MakeMap(value.Type())

	if len(input) > 0 {
		input = strings.Trim(input, "[]")

		// An empty string would cause an empty map
		r := csv.NewReader(strings.NewReader(input))
		ss, err := r.Read()
		if err != nil {
			return err
		}

		for _, pair := range ss {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return errors.Errorf("%s must be formatted as key=value", pair)
			}

			m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(kv[1]))
		}
	}

	value.Set(m)
	return nil
}

func base64ByteSliceFieldSetter(value reflect.Value, input string) error {
	if len(input) > 0 {
		b, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(b))
	}
	return nil
}

func hexByteSliceFieldSetter(value reflect.Value, input string) error {
	if len(input) > 0 {
		b, err := hex.DecodeString(input)
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(b))
	}
	return nil
}

func sliceFieldSetter(value reflect.Value, input string) error {
	input = strings.Trim(input, "[]")

	if len(input) > 0 {
		r := csv.NewReader(strings.NewReader(input))
		ss, err := r.Read()
		if err != nil {
			return err
		}

		s := reflect.MakeSlice(value.Type(), len(ss), len(ss))

		for n, pair := range ss {
			s.Index(n).Set(reflect.ValueOf(pair).Convert(value.Type().Elem()))
		}

		value.Set(s)
	}

	return nil
}

func stringFieldSetter(value reflect.Value, input string) error {
	value.SetString(input)
	return nil
}

func timeDurationFieldSetter(value reflect.Value, input string) error {
	if len(input) > 0 {
		d, err := time.ParseDuration(input)
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(d))
	}

	return nil
}

func pointerFieldSetter(x func(value reflect.Value, input string) error) func(value reflect.Value, input string) error {
	return func(value reflect.Value, input string) error {
		if value.Kind() != reflect.Ptr {
			return x(value, input)
		}

		// Create a new instance of the specified type
		v := reflect.New(value.Type().Elem())

		// Set into that new instance
		err := x(v.Elem(), input)
		if err != nil {
			return err
		}

		// Set that instance to the pointer
		value.Set(v)

		return nil
	}
}

func init() {
	knownSetters = make(map[string]fieldSetter)
	RegisterConfigParser("time.Duration", timeDurationFieldSetter)
}
