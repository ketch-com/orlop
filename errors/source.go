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

type sourcer struct {
	error
	source string
}

func (s sourcer) Unwrap() error {
	return s.error
}

func (s sourcer) Source() string {
	return s.source
}

// WithSource adds a Sourcer to err's error chain.
// Unlike pkg/errors, WithSource will wrap nil error.
func WithSource(err error, source string) error {
	if err == nil {
		err = New("Source<" + source + ">")
	}
	return sourcer{err, source}
}

// Source returns the source associated with an error.
// If err is nil, it returns "".
func Source(err error) string {
	var um internal.Source

	if err != nil {
		if As(err, &um) {
			return um.Source()
		}
	}

	return ""
}
