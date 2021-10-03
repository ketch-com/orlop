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

package grpc

import (
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.ketch.com/lib/orlop/v2/errors"
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

	if rv.Kind() != reflect.Ptr {
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

		if rv.Kind() != reflect.Ptr {
			panic("binary: v must be a pointer")
		}

		if rv.IsNil() {
			return errors.New("binary: v must not be nil")
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
