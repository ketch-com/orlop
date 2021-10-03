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
	"fmt"
	"strings"
)

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
