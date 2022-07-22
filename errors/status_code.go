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
	"net/http"
)

type statusCoder struct {
	error
	code int
}

func (sc statusCoder) Cause() error {
	return sc.error
}

func (sc statusCoder) Unwrap() error {
	return sc.error
}

func (sc statusCoder) StatusCode() int {
	return sc.code
}

func (sc statusCoder) Timeout() bool {
	var temp internal.Timeout
	if As(sc.error, &temp) {
		return temp.Timeout()
	}

	return sc.code == http.StatusRequestTimeout
}

func (sc statusCoder) Temporary() bool {
	var temp internal.Temporary
	if As(sc.error, &temp) {
		return temp.Temporary()
	}

	return sc.code == http.StatusInternalServerError || sc.code == http.StatusServiceUnavailable || sc.code == http.StatusRequestTimeout
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
// If err is nil, it returns 200 http.StatusOK.
func StatusCode(err error) (code int) {
	var sc internal.StatusCode
	var ec internal.ErrorCode

	if err == nil {
		return http.StatusOK
	}
	if As(err, &sc) {
		return sc.StatusCode()
	}
	if As(err, &ec) {
		switch ec.ErrorCode() {
		case ECONFLICT:
			return http.StatusConflict

		case EUNAVAILABLE:
			return http.StatusServiceUnavailable

		case EINVALID:
			return http.StatusBadRequest

		case ENOTFOUND:
			return http.StatusNotFound

		case ETIMEOUT:
			return http.StatusRequestTimeout

		case EFORBIDDEN:
			return http.StatusForbidden
		}
	}

	return http.StatusInternalServerError
}
