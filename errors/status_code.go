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
	"net/http"
)

type statusCoder struct {
	error
	code int
}

func (sc statusCoder) Unwrap() error {
	return sc.error
}

func (sc statusCoder) Error() string {
	return fmt.Sprintf("[%d] %v", sc.code, sc.error)
}

func (sc statusCoder) StatusCode() int {
	return sc.code
}

// WithStatusCode adds a StatusCoder to err's error chain.
// Unlike pkg/errors, WithStatusCode will wrap nil error.
func WithStatusCode(err error, code int) error {
	if err == nil {
		err = New(http.StatusText(code))
	}
	return statusCoder{err, code}
}

// StatusCode returns the status code associated with an error.
// If no status code is found, it returns 500 http.StatusInternalServerError.
// As a special case, it checks for Timeout() and Temporary() errors and returns
// 504 http.StatusGatewayTimeout and 503 http.StatusServiceUnavailable
// respectively.
// If err is nil, it returns 200 http.StatusOK.
func StatusCode(err error) (code int) {
	if err == nil {
		return http.StatusOK
	}

	var sc interface {
		error
		StatusCode() int
	}

	if As(err, &sc) {
		return sc.StatusCode()
	}

	var timeouter interface {
		error
		Timeout() bool
	}
	if As(err, &timeouter) && timeouter.Timeout() {
		return http.StatusGatewayTimeout
	}

	var temper interface {
		error
		Temporary() bool
	}
	if As(err, &temper) && temper.Temporary() {
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}
