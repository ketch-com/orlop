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
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/fx"

	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/service"
)

type Config interface {
	Options() fx.Option
}

type value struct {
	isPopulated bool
	value       any
}

type providerImpl struct {
	configs map[string]value
	environ env.Environ
	prefix  service.Name
}

func New(p Params) Provider {
	configs := make(map[string]value)
	for _, c := range p.Defs {
		configs[strings.ToLower(c.Name)] = value{
			isPopulated: false,
			value:       c.Config,
		}
	}

	return &providerImpl{
		configs: configs,
		environ: p.Environ,
		prefix:  p.Prefix,
	}
}

func (s *providerImpl) Get(ctx context.Context, service string) (any, error) {
	serviceName := strings.ToLower(service)
	if cfg, ok := s.configs[serviceName]; ok {
		if cfg.isPopulated {
			return cfg.value, nil
		}

		err := s.load(ctx, serviceName, cfg.value)
		if err != nil {
			return cfg.value, err
		}
		return cfg.value, nil
	}

	return nil, errors.Errorf("%s config not found", service)
}

func (s *providerImpl) List(_ context.Context) ([]string, error) {
	var vars []string
	for k, v := range s.configs {
		vs, err := s.getVariablesFromConfig(k, v.value)
		if err != nil {
			return nil, err
		}

		vars = append(vars, vs...)
	}

	return vars, nil
}

func (s *providerImpl) load(_ context.Context, key string, value any) error {
	fields, err := reflectStruct([]string{}, value)
	if err != nil {
		return err
	}

	key = strings.TrimSpace(key)
	for name, field := range fields {
		keyName := name
		if len(key) > 0 {
			keyName = strings.Join([]string{key, name}, "_")
		}

		if v := s.environ.Getenv(keyName); len(v) > 0 {
			if err = field.set(field.v, v); err != nil {
				return errors.Wrapf(err, "failed to set field '%s' with value '%s'", name, v)
			}
		} else if field.tag.DefaultValue != nil {
			if err = field.set(field.v, *field.tag.DefaultValue); err != nil {
				return errors.Wrapf(err, "failed to set field '%s' with value '%s'", name, v)
			}
		} else if field.tag.Required {
			return errors.Errorf("%s required", name)
		}
	}

	return nil
}

// getVariablesFromConfig returns the environment variables from the given config object
func (s *providerImpl) getVariablesFromConfig(service string, cfg any) ([]string, error) {
	var vars []string

	prefix := []string{string(s.prefix)}
	if len(strings.TrimSpace(service)) > 0 {
		prefix = append(prefix, service)
	}

	fields, err := reflectStruct(prefix, cfg)
	if err != nil {
		return nil, err
	}

	for name, field := range fields {
		var value string
		if field.tag.DefaultValue != nil {
			value = *field.tag.DefaultValue
		}

		if len(value) == 0 {
			switch field.v.Kind() {
			case reflect.Bool:
				value = "false # bool"

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value = "0 # int"

			case reflect.Float32, reflect.Float64:
				value = "0.0 # float"

			case reflect.Map:
				value = "# [k=v, k=v, k=v]"

			case reflect.Slice:
				if field.v.Type().Elem().Kind() == reflect.Int8 {
					value = "# bytes"
				} else {
					value = "# [v1, v2, v3]"
				}

			case reflect.String:
				value = "# string"

			}
		}

		vars = append(vars, fmt.Sprintf("%s=%s", name, value))
	}

	return vars, nil
}

func toScreamingDelimited(s string, delimiter uint8, ignore uint8, screaming bool) string {
	s = strings.TrimSpace(s)
	n := strings.Builder{}
	n.Grow(len(s) + 2) // nominal 2 bytes of extra space for inserted delimiters
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if vIsLow && screaming {
			v += 'A'
			v -= 'a'
		} else if vIsCap && !screaming {
			v += 'a'
			v -= 'A'
		}

		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		if i+1 < len(s) {
			next := s[i+1]
			vIsNum := v >= '0' && v <= '9'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			nextIsNum := next >= '0' && next <= '9'
			// add underscore if next letter case type is changed
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) || (vIsNum && (nextIsCap || nextIsLow)) {
				if prevIgnore := ignore > 0 && i > 0 && s[i-1] == ignore; !prevIgnore {
					if vIsCap && nextIsLow {
						if prevIsCap := i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z'; prevIsCap {
							n.WriteByte(delimiter)
						}
					}
					n.WriteByte(v)
					if vIsLow || vIsNum || nextIsNum {
						n.WriteByte(delimiter)
					}
					continue
				}
			}
		}

		if (v == ' ' || v == '_' || v == '-') && uint8(v) != ignore {
			// replace space/underscore/hyphen with delimiter
			n.WriteByte(delimiter)
		} else {
			n.WriteByte(v)
		}
	}

	return n.String()
}
