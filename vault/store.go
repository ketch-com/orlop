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

package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/parameter"
	"path"
	"reflect"
	"strings"
)

type client struct {
	params Params
	client *vault.Client
}

// ReadObject returns a secret at the given path
func (c client) ReadObject(ctx context.Context, p string, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		panic("out must be a pointer")
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	sec, err := c.Read(ctx, p)
	if err != nil {
		return err
	}

	if sec == nil {
		return parameter.ErrNotFound
	}

	t := v.Type()
	for n := 0; n < t.NumField(); n++ {
		ft := t.Field(n)

		name := ft.Name

		parts := strings.Split(ft.Tag.Get("json"), ",")
		if len(parts) > 0 && len(parts[0]) > 0 {
			name = parts[0]
		}

		if val, ok := sec[name]; ok {
			f := v.Field(n)
			switch f.Kind() {
			case reflect.String:
				if s, ok := val.(string); ok {
					f.SetString(s)
				} else {
					f.SetString(fmt.Sprintf("%v", val))
				}

			case reflect.Bool:
				if b, ok := val.(bool); ok {
					f.SetBool(b)
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if j, ok := val.(json.Number); ok {
					if i, err := j.Int64(); err == nil {
						f.SetInt(i)
					}
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if j, ok := val.(json.Number); ok {
					if i, err := j.Int64(); err == nil {
						f.SetUint(uint64(i))
					}
				}

			case reflect.Float32, reflect.Float64:
				if j, ok := val.(json.Number); ok {
					if i, err := j.Float64(); err == nil {
						f.SetFloat(i)
					}
				}

			case reflect.Slice:
				if j, ok := val.(string); ok {
					b, err := base64.StdEncoding.DecodeString(j)
					if err != nil {
						return err
					}

					f.SetBytes(b)
				} else if j, ok := val.([]interface{}); ok {
					sl := reflect.MakeSlice(f.Type(), 0, len(j))
					for _, sv := range j {
						sl = reflect.Append(sl, reflect.ValueOf(sv))
					}
					f.Set(sl)
				}

			default:
				panic(fmt.Sprintf("%s not supported", f.Type()))
			}
		}
	}

	return nil
}

// Read returns a secret at the given path
func (c client) Read(ctx context.Context, p string) (map[string]interface{}, error) {
	_, span := c.params.Tracer.Start(ctx, "Read")
	defer span.End()

	keyPath := c.prefix(p)

	sec, err := c.client.Logical().Read(keyPath)
	if err != nil {
		c.params.Logger.WithField("vault.path", keyPath).WithError(err).Trace("read failed")
		span.RecordError(err)
		return nil, errors.Wrap(err, "failed to read from Vault")
	}

	return sec.Data, nil
}

// WriteObject writes secret data at the given path from an object
func (c client) WriteObject(ctx context.Context, p string, in interface{}) error {
	m := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Struct {
		panic("in must be a pointer or struct")
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("in must be a pointer or struct")
	}

	t := v.Type()
	for n := 0; n < t.NumField(); n++ {
		ft := t.Field(n)

		name := ft.Name

		parts := strings.Split(ft.Tag.Get("json"), ",")
		if len(parts) > 0 && len(parts[0]) > 0 {
			name = parts[0]
		}

		f := v.Field(n)
		switch f.Kind() {
		case reflect.String:
			m[name] = f.String()

		case reflect.Bool:
			m[name] = f.Bool()

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[name] = f.Int()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			m[name] = f.Uint()

		case reflect.Float32, reflect.Float64:
			m[name] = f.Float()

		case reflect.Slice:
			ftek := f.Type().Elem().Kind()
			if ftek == reflect.Uint8 {
				m[name] = base64.StdEncoding.EncodeToString(f.Bytes())
			} else {
				var sl []interface{}

				for n := 0; n < f.Len(); n++ {
					fn := f.Index(n)
					var fv interface{}

					switch ftek {
					case reflect.String:
						fv = fn.String()

					case reflect.Bool:
						fv = fn.Bool()

					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						fv = fn.Int()

					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						fv = fn.Uint()

					case reflect.Float32, reflect.Float64:
						fv = fn.Float()
					}

					sl = append(sl, fv)
				}

				m[name] = sl
			}

		default:
			panic(fmt.Sprintf("%s not supported", f.Type()))
		}
	}

	if _, err := c.Write(ctx, p, m); err != nil {
		return err
	}

	return nil
}

// Write secret data at the given path
func (c client) Write(ctx context.Context, p string, data map[string]interface{}) (map[string]interface{}, error) {
	_, span := c.params.Tracer.Start(ctx, "Write")
	defer span.End()

	keyPath := c.prefix(p)

	sec, err := c.client.Logical().Write(keyPath, data)
	if err != nil {
		c.params.Logger.WithField("vault.path", keyPath).WithError(err).Trace("write failed")
		span.RecordError(err)
		return nil, errors.Wrap(err, "failed to write to Vault")
	}

	return sec.Data, nil
}

// Delete a secret at the given path
func (c client) Delete(ctx context.Context, p string) error {
	_, span := c.params.Tracer.Start(ctx, "Delete")
	defer span.End()

	keyPath := c.prefix(p)

	_, err := c.client.Logical().Delete(keyPath)
	if err != nil {
		c.params.Logger.WithField("vault.path", keyPath).WithError(err).Trace("delete failed")
		span.RecordError(err)
		return errors.Wrap(err, "failed to delete from Vault")
	}

	return nil
}

// List returns keys available at the given path p
func (c client) List(ctx context.Context, p string) ([]string, error) {
	_, span := c.params.Tracer.Start(ctx, "List")
	defer span.End()

	var keys []string

	keyPath := c.prefix(p)

	sec, err := c.client.Logical().List(keyPath)
	if err != nil {
		c.params.Logger.WithField("vault.path", keyPath).WithError(err).Trace("list failed")
		span.RecordError(err)
		return nil, errors.Wrap(err, "failed to list from Vault")
	}

	if sec == nil || sec.Data == nil || sec.Data["keys"] == nil {
		return nil, nil
	}

	for _, key := range sec.Data["keys"].([]interface{}) {
		keys = append(keys, key.(string))
	}

	return keys, nil
}

func (c client) GetEnabled() bool {
	return true
}

func (c client) prefix(p string) string {
	return path.Join(c.params.Config.Prefix, p)
}
