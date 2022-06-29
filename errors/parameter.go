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

package errors

import (
	"go.ketch.com/lib/orlop/v2/errors/internal"
)

type parameter struct {
	error
	parameter string
}

func (s parameter) Unwrap() error {
	return s.error
}

func (s parameter) Parameter() string {
	return s.parameter
}

// WithParameter adds a parameter to err's error chain.
// Unlike pkg/errors, WithParameter will wrap nil error.
func WithParameter(err error, param string) error {
	if err == nil {
		err = New("Parameter<" + param + ">")
	}
	return parameter{err, param}
}

// Parameter returns the parameter associated with an error.
// If err is nil, it returns "".
func Parameter(err error) string {
	var um internal.Parameter

	if err != nil {
		if As(err, &um) {
			return um.Parameter()
		}
	}

	return ""
}
