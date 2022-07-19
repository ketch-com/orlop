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
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"go.ketch.com/lib/orlop/v2/errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Unmarshal reads configuration into the cfg object
//
// deprecated: Unmarshal not supported and will be removed soon
func Unmarshal(prefix string, cfg any) error {
	return UnmarshalFromEnv(prefix, os.Environ(), cfg)
}

// UnmarshalFromEnv reads configuration into the cfg object from the env vars
//
// deprecated: UnmarshalFromEnv not supported and will be removed soon
func UnmarshalFromEnv(prefix string, vars []string, cfg any) error {
	prefix = toScreamingDelimited(prefix, '_', 0, true)

	env := make(map[string]string)
	for _, v := range vars {
		parts := strings.SplitN(v, "=", 2)
		if strings.HasPrefix(parts[0], prefix) {
			env[parts[0]] = parts[1]
		}
	}

	fields, err := reflectStruct([]string{prefix}, cfg)
	if err != nil {
		return err
	}

	for name, field := range fields {
		if v, ok := env[name]; ok {
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

type configTag struct {
	Name         *string
	Encoding     *string
	DefaultValue *string
	Required     bool
}

func (t configTag) String() string {
	var s []string
	if t.Name != nil {
		s = append(s, fmt.Sprintf("name=%s", *t.Name))
	}
	if t.Encoding != nil {
		s = append(s, fmt.Sprintf("encoding=%s", *t.Encoding))
	}
	if t.DefaultValue != nil {
		s = append(s, fmt.Sprintf("defaultValue=%s", *t.DefaultValue))
	}
	if t.Required {
		s = append(s, "required")
	} else {
		s = append(s, "optional")
	}

	return strings.Join(s, ", ")
}

func parseConfigTag(tag string) *configTag {
	t := &configTag{}

	if len(tag) == 0 {
		return t
	}

	parts := strings.Split(tag, ",")

	t.Name = &parts[0]

	for _, part := range parts[1:] {
		elemParts := strings.Split(part, "=")
		switch elemParts[0] {
		case "default":
			t.DefaultValue = &elemParts[1]

		case "required":
			t.Required = true

		case "encoding":
			switch elemParts[1] {
			case "base64", "hex":
				t.Encoding = &elemParts[1]

			default:
				panic("config: unsupported encoding in " + tag)
			}
		}
	}

	return t
}

type fieldSetter func(value reflect.Value, input string) error

var knownSetters map[string]fieldSetter

// RegisterConfigParser registers a config parser
//
// deprecated: RegisterConfigParser not supported and will be removed soon
func RegisterConfigParser(typeName string, parser func(value reflect.Value, input string) error) {
	knownSetters[typeName] = parser
}

type configField struct {
	tag *configTag
	v   reflect.Value
	set fieldSetter
}

func reflectStruct(prefix []string, i any) (map[string]*configField, error) {
	r := make(map[string]*configField)

	err := reflectStructValue(prefix, r, reflect.ValueOf(i))
	if err != nil {
		return nil, err
	}

	return r, nil
}

func reflectStructValue(prefix []string, r map[string]*configField, v reflect.Value) error {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	t := v.Type()

	for n := 0; n < t.NumField(); n++ {
		f := v.Field(n)
		ft := f.Type()
		fld := t.Field(n)
		tag := parseConfigTag(fld.Tag.Get("config"))

		if tag.Name == nil || len(*tag.Name) == 0 && ft.Kind() != reflect.Struct {
			nm := fld.Name
			tag.Name = &nm
		}

		if *tag.Name == "-" {
			continue
		}

		replacer := strings.NewReplacer("-", "_", ".", "_")
		key := toScreamingDelimited(replacer.Replace(strings.Join(append(prefix, *tag.Name), "_")), '_', 0, true)

		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		setter := knownSetters[ft.String()]
		if f.CanAddr() {
			if setter == nil {
				m := f.Addr().MethodByName("UnmarshalText")
				if m.IsValid() {
					setter = unmarshalTextSetter
				}
			}
			if setter == nil {
				m := f.Addr().MethodByName("UnmarshalJSON")
				if m.IsValid() {
					setter = unmarshalJSONSetter
				}
			}
		}

		if setter == nil {
			switch ft.Kind() {
			case reflect.Bool:
				setter = boolFieldSetter

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				setter = intFieldSetter

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				setter = uintFieldSetter

			case reflect.Float32, reflect.Float64:
				setter = floatFieldSetter

			case reflect.Map:
				setter = mapFieldSetter

			case reflect.Slice:
				if ft.Elem().Kind() == reflect.Uint8 {
					if tag.Encoding != nil && *tag.Encoding == "base64" {
						setter = base64ByteSliceFieldSetter
					} else {
						setter = hexByteSliceFieldSetter
					}
				} else {
					setter = sliceFieldSetter
				}

			case reflect.Struct:
				p := prefix
				if len(*tag.Name) > 0 {
					p = append(prefix, *tag.Name)
				}

				err := reflectStructValue(p, r, f)
				if err != nil {
					return err
				}

			case reflect.String:
				setter = stringFieldSetter
			}
		}

		if setter != nil {
			if f.Type().Kind() == reflect.Ptr {
				setter = pointerFieldSetter(setter)
			}

			r[key] = &configField{
				tag: tag,
				v:   f,
				set: setter,
			}
		}
	}

	return nil
}

// GetVariablesFromConfig returns the environment variables from the given config object
//
// deprecated: GetVariablesFromConfig not supported and will be removed soon
func GetVariablesFromConfig(prefix string, cfg any) ([]string, error) {
	var vars []string

	fields, err := reflectStruct([]string{prefix}, cfg)
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

// LoadEnvironment loads the environment variables from the specified files and from the standard locations
func LoadEnvironment(env Environment, files ...string) {
	var envFiles []string
	envFiles = append(envFiles, files...)
	if env.IsLocal() {
		envFiles = append(envFiles, ".env.local")
	} else {
		envFiles = append(envFiles, ".env."+env.String())
	}
	envFiles = append(envFiles, ".env")

	for _, file := range envFiles {
		if _, err := os.Stat(file); err == nil {
			_ = godotenv.Load(file)
		}
	}
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

func init() {
	knownSetters = make(map[string]fieldSetter)
	RegisterConfigParser("time.Duration", timeDurationFieldSetter)
}
