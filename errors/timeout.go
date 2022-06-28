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
	"net/http"
	"time"
)

type timeouter struct {
	error
	timeout time.Duration
}

func (te timeouter) Error() string {
	if te.timeout > time.Duration(0) {
		return fmt.Sprintf("timed out after %v - %v", te.timeout, te.error)
	}

	return fmt.Sprintf("timed out - %v", te.error)
}

func (te timeouter) Unwrap() error {
	return te.error
}

func (te timeouter) UserMessage() string {
	return "operation timed out"
}

func (te timeouter) Timeout() time.Duration {
	return te.timeout
}

// Timeout returns a timeout with ETIMEOUT, 504 Gateway Timeout and a user message "operation timed out"
func Timeout(err error, timeout ...time.Duration) error {
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &timeouter{err, t}
}

// IsTimeout returns true if the error is a Timeout error
func IsTimeout(err error) bool {
	var timeouter internal.Timeout
	var sc internal.StatusCode
	var ec internal.ErrorCode

	if As(err, &timeouter) {
		return true
	}
	if As(err, &sc) && (sc.StatusCode() == http.StatusGatewayTimeout || sc.StatusCode() == http.StatusGatewayTimeout) {
		return true
	}
	if As(err, &ec) && ec.ErrorCode() == ETIMEOUT {
		return true
	}

	return false
}
