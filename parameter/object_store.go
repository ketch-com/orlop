package parameter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ObjectStore is an extension to read/write objects to parameter store
type ObjectStore interface {
	Store

	ReadObject(ctx context.Context, p string, out interface{}) error
	WriteObject(ctx context.Context, p string, in interface{}) error
}

type objectStore struct {
	store Store
}

func NewObjectStore(store Store) ObjectStore {
	return &objectStore{
		store: store,
	}
}

func (c *objectStore) List(ctx context.Context, p string) ([]string, error) {
	return c.store.List(ctx, p)
}

func (c *objectStore) Read(ctx context.Context, p string) (map[string]interface{}, error) {
	return c.store.Read(ctx, p)
}

func (c *objectStore) Write(ctx context.Context, p string, data map[string]interface{}) error {
	return c.store.Write(ctx, p, data)
}

func (c *objectStore) Delete(ctx context.Context, p string) error {
	return c.store.Delete(ctx, p)
}

// ReadObject returns a secret at the given path
func (c *objectStore) ReadObject(ctx context.Context, p string, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		panic("out must be a pointer")
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	sec, err := c.store.Read(ctx, p)
	if err != nil {
		return err
	}

	if sec == nil {
		return ErrNotFound
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

// WriteObject writes secret data at the given path from an object
func (c *objectStore) WriteObject(ctx context.Context, p string, in interface{}) error {
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

	if err := c.store.Write(ctx, p, m); err != nil {
		return err
	}

	return nil
}
