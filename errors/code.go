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
	"fmt"
)

type coder struct {
	error
	code ErrorCode
}

func (sc coder) Unwrap() error {
	return sc.error
}

func (sc coder) Error() string {
	return fmt.Sprintf("[%s] %v", sc.code, sc.error)
}

func (sc coder) ErrorCode() ErrorCode {
	return sc.code
}

// WithCode adds a Coder to err's error chain.
// Unlike pkg/errors, WithCode will wrap nil error.
func WithCode(err error, code ErrorCode) error {
	if err == nil {
		err = New(string(code))
	}
	return coder{err, code}
}

// Code returns the code associated with an error.
// If no code is found, it returns EINTERNAL.
// As a special case, it checks for Timeout() and Temporary() errors and returns
// ETIMEOUT and EUNAVAILABLE
// respectively.
// If err is nil, it returns empty string
func Code(err error) (code ErrorCode) {
	if err == nil {
		return ""
	}

	var sc interface {
		error
		ErrorCode() ErrorCode
	}

	if As(err, &sc) {
		return sc.ErrorCode()
	}

	var timeouter interface {
		error
		Timeout() bool
	}
	if As(err, &timeouter) && timeouter.Timeout() {
		return ETIMEOUT
	}

	var temper interface {
		error
		Temporary() bool
	}
	if As(err, &temper) && temper.Temporary() {
		return EUNAVAILABLE
	}

	return EINTERNAL
}
