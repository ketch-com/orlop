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
	"go.ketch.com/lib/orlop/v2/errors/internal"
)

type coder struct {
	error
	code string
}

func (sc coder) Unwrap() error {
	return sc.error
}

func (sc coder) Error() string {
	return fmt.Sprintf("[%s] %v", sc.code, sc.error)
}

func (sc coder) ErrorCode() string {
	return sc.code
}

// WithCode adds a Coder to err's error chain.
// Unlike pkg/errors, WithCode will wrap nil error.
func WithCode(err error, code string) error {
	if err == nil {
		err = New(code)
	}
	return coder{err, code}
}

// Code returns the code associated with an error.
// If no code is found, it returns EINTERNAL.
// As a special case, it checks for Timeout() and Temporary() errors and returns
// ETIMEOUT and EUNAVAILABLE
// respectively.
// If err is nil, it returns empty string
func Code(err error) (code string) {
	var ec internal.ErrorCode

	if err == nil {
		return ""
	}
	if As(err, &ec) {
		return ec.ErrorCode()
	}
	if IsTimeout(err) {
		return ETIMEOUT
	}
	if IsTemporary(err) {
		return EUNAVAILABLE
	}
	if IsNotFound(err) {
		return ENOTFOUND
	}

	return EINTERNAL
}
