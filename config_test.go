// Copyright (c) 2020 Ketch, Inc.
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
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type EmbeddedConfig struct {
	Embedded bool
}

type CustomConfig struct {
	Custom time.Duration
}

//func (e *CustomConfig) UnmarshalText(text []byte) error {
//	fmt.Printf("Called UnmarshalText with %v\n", text)
//	return nil
//}

func (e *CustomConfig) UnmarshalJSON(text []byte) error {
	fmt.Printf("Called UnmarshalJSON with %v\n", text)
	return nil
}

type LargerConfig struct {
	CustomConfig CustomConfig
	L0Base       int32
}

type TestConfig struct {
	Embedded      EmbeddedConfig `config:","`
	CustomConfig  CustomConfig   `config:"unmarshaller"`
	WithDefault   string         `config:"def,default=/pki/issue"`
	Required      string         `config:"req,required"`
	SomeSlice     []string       `config:"sliced"`
	CustomParser  time.Duration  `config:"custom,default=12345s"`
	Map           map[string]string
	HexEncoded    []byte
	Base64Encoded []byte `config:",encoding=base64"`
	Ptr           *int32
}

func (e *TestConfig) UnmarshalText(text []byte) error {
	fmt.Printf("Called UnmarshalText with %v\n", text)
	return nil
}

func TestUnmarshalStruct(t *testing.T) {
	err := Unmarshal("wheelhouse", &LargerConfig{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVars(t *testing.T) {
	vars, err := GetVariablesFromConfig("wheelhouse", &TestConfig{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(strings.Join(vars, "\n"))
}

func TestUnmarshal(t *testing.T) {
	var c TestConfig

	err := UnmarshalFromEnv("wheelhouse", []string{
		"WHEELHOUSE_EMBEDDED=true",
		"WHEELHOUSE_REQ=imhere",
		"WHEELHOUSE_SLICED=\"a\",\"b\",\"c\"",
		"WHEELHOUSE_CUSTOM=1m",
		"WHEELHOUSE_MAP=[\"a=b\",\"c=d\"]",
		"WHEELHOUSE_HEX_ENCODED=0102030405060708090A0B0C0D0E0F",
		"WHEELHOUSE_BASE_64_ENCODED=AQIDBAUGBwgJCgsMDQ4P",
		"WHEELHOUSE_PTR=123",
	}, &c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)

	if !c.Embedded.Embedded {
		t.Fail()
	}

	if c.WithDefault != "/pki/issue" {
		t.Fail()
	}

	if c.Required != "imhere" {
		t.Fail()
	}

	if c.SomeSlice == nil || !reflect.DeepEqual(c.SomeSlice, []string{"a", "b", "c"}) {
		t.Fail()
	}

	if c.CustomParser != time.Minute {
		t.Fail()
	}

	if !reflect.DeepEqual(c.Map, map[string]string{"a": "b", "c": "d"}) {
		t.Fail()
	}

	if !reflect.DeepEqual(c.HexEncoded, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}) {
		t.Fail()
	}

	if c.Ptr == nil || *c.Ptr != 123 {
		t.Fail()
	}
}

func st(s string) *string {
	return &s
}

func TestParseConfigTag(t *testing.T) {
	for _, f := range []struct {
		S string
		T *configTag
	}{
		{
			",",
			&configTag{
				Name: st(""),
			},
		},
		{
			"nm",
			&configTag{
				Name: st("nm"),
			},
		},
		{
			",required",
			&configTag{
				Name:     st(""),
				Required: true,
			},
		},
		{
			",default=foo",
			&configTag{
				Name:         st(""),
				DefaultValue: st("foo"),
			},
		},
		{
			",encoding=hex",
			&configTag{
				Name:     st(""),
				Encoding: st("hex"),
			},
		},
		{
			",encoding=hex",
			&configTag{
				Name:     st(""),
				Encoding: st("hex"),
			},
		},
		{
			"name,default=foo,encoding=hex,required",
			&configTag{
				Name:         st("name"),
				DefaultValue: st("foo"),
				Encoding:     st("hex"),
				Required:     true,
			},
		},
	} {
		if f.T.String() != parseConfigTag(f.S).String() {
			t.Fatal(f.S, f.T.String(), parseConfigTag(f.S).String())
		}
	}
}
