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

package internal

import (
	"errors"
	"google.golang.org/grpc/status"
)

// Application error codes.
const (
	ECONFLICT      = "conflict"      // action cannot be performed
	EINTERNAL      = "internal"      // internal error
	EUNAVAILABLE   = "unavailable"   // service is unavailable
	EINVALID       = "invalid"       // validation failed
	ENOTFOUND      = "not_found"     // entity does not exist
	ETIMEOUT       = "timeout"       // operation timed out
	ECANCELED      = "canceled"      // operation canceled
	EFORBIDDEN     = "forbidden"     // operation is not authorized
	ECONFIGURATION = "configuration" // configuration error
	EUNIMPLEMENTED = "unimplemented" // unimplemented error
)

func NewStandardError(err error, code string) error {
	return &StandardError{err, code}
}

type StandardError struct {
	error
	code string
}

func (c StandardError) Cause() error {
	return c.error
}

func (c StandardError) Unwrap() error {
	return c.error
}

func (c StandardError) ErrorCode() string {
	return c.code
}

func (c StandardError) StatusCode() int {
	return StandardToHTTP[c.code]
}

func (c StandardError) Status() int {
	return StandardToHTTP[c.code]
}

func (c StandardError) GRPCStatus() *status.Status {
	return status.New(StandardToGrpc[c.code], c.Error())
}

func (c StandardError) Timeout() bool {
	var to Timeout
	if errors.As(c.error, &to) {
		return to.Timeout()
	}

	return c.code == ETIMEOUT
}

func (c StandardError) Temporary() bool {
	var temp Temporary
	if errors.As(c.error, &temp) {
		return temp.Temporary()
	}

	return c.code == EINTERNAL || c.code == EUNAVAILABLE || c.code == ETIMEOUT
}
