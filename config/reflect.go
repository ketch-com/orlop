package config

import (
	"go.ketch.com/lib/orlop/errors"
	"reflect"
	"strings"
)

type configField struct {
	tag *configTag
	v   reflect.Value
	set fieldSetter
}

func reflectStruct(prefix []string, i interface{}) (map[string]*configField, error) {
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
