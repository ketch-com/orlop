package orlop

import (
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"io"
	"io/ioutil"
	"reflect"
)

// BinaryMarshaler marshals the given object as a binary object
type BinaryMarshaler struct{}

// Marshal marshals "v" into byte sequence.
func (BinaryMarshaler) Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice || rv.Elem().Kind() != reflect.Uint8 {
		return json.Marshal(v)
	}

	return v.([]byte), nil
}

// Unmarshal unmarshals "data" into "v". "v" must be a pointer value.
func (BinaryMarshaler) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)

	for rv.Kind() != reflect.Ptr {
		panic("binary: v must be a pointer")
	}

	if rv.IsNil() {
		panic("binary: v must not be nil")
	}

	rv.Elem().SetBytes(data)

	return nil
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func (BinaryMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(v interface{}) error {
		rv := reflect.ValueOf(v)

		for rv.Kind() != reflect.Ptr {
			panic("binary: v must be a pointer")
		}

		if rv.IsNil() {
			return fmt.Errorf("binary: v must not be nil")
		}

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		rv.Elem().SetBytes(data)

		return nil
	})
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
func (BinaryMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(v interface{}) error {
		rv := reflect.ValueOf(v)

		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return nil
			}
			rv = rv.Elem()
		}

		if rv.Kind() != reflect.Slice || rv.Elem().Kind() != reflect.Uint8 {
			return json.NewEncoder(w).Encode(v)
		}

		_, err := w.Write(v.([]byte))

		return err
	})
}

// ContentType returns the Content-Type which this marshaler is responsible for.
func (BinaryMarshaler) ContentType() string {
	return "application/octet-stream"
}
