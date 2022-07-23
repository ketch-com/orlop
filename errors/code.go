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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// GRPCStatus() *status.Status
// 	if errors.Is(err, context.DeadlineExceeded) {
//		return New(codes.DeadlineExceeded, err.Error())
//	}
//	if errors.Is(err, context.Canceled) {
//		return New(codes.Canceled, err.Error())
//	}

type coder struct {
	error
	code string
}

func (c coder) Cause() error {
	return c.error
}

func (c coder) Unwrap() error {
	return c.error
}

func (c coder) ErrorCode() string {
	return c.code
}

func (c coder) Timeout() bool {
	var to internal.Timeout
	if As(c.error, &to) {
		return to.Timeout()
	}

	return c.code == ETIMEOUT
}

func (c coder) Temporary() bool {
	var temp internal.Temporary
	if As(c.error, &temp) {
		return temp.Temporary()
	}

	return c.code == EINTERNAL || c.code == EUNAVAILABLE || c.code == ETIMEOUT
}

func (c coder) GRPCStatus() *status.Status {
	var code codes.Code

	switch c.code {
	case ECONFLICT:
		code = codes.Aborted

	case ECANCELED:
		code = codes.Canceled

	case EINTERNAL:
		code = codes.Internal

	case EUNAVAILABLE:
		code = codes.Unavailable

	case EINVALID:
		code = codes.InvalidArgument

	case ENOTFOUND:
		code = codes.NotFound

	case ETIMEOUT:
		code = codes.DeadlineExceeded

	default:
		code = codes.Unknown
	}

	return status.New(code, c.Error())
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
// If err is nil, it returns empty string
func Code(err error) string {
	var sc internal.StatusCode
	var ec internal.ErrorCode
	var gc internal.GRPCStatus

	if err == nil {
		return ""
	}
	if As(err, &ec) {
		return ec.ErrorCode()
	}
	if As(err, &sc) {
		switch sc.StatusCode() {
		case http.StatusBadRequest:
			return EINVALID

		case http.StatusUnauthorized:
			return EFORBIDDEN

		case http.StatusPaymentRequired:
			return EFORBIDDEN

		case http.StatusForbidden:
			return EFORBIDDEN

		case http.StatusNotFound:
			return ENOTFOUND

		case http.StatusMethodNotAllowed:
			return EINVALID

		case http.StatusNotAcceptable:
			return EINVALID

		case http.StatusProxyAuthRequired:
			return EFORBIDDEN

		case http.StatusRequestTimeout:
			return ETIMEOUT

		case http.StatusConflict:
			return ECONFLICT

		case http.StatusGone:
			return ENOTFOUND

		case http.StatusLengthRequired:
			return EINVALID

		case http.StatusPreconditionFailed:
			return EINVALID

		case http.StatusRequestEntityTooLarge:
			return EINVALID

		case http.StatusRequestURITooLong:
			return EINVALID

		case http.StatusUnsupportedMediaType:
			return EINVALID

		case http.StatusRequestedRangeNotSatisfiable:
			return EINVALID

		case http.StatusExpectationFailed:
			return EINVALID

		case http.StatusTeapot:
			return EINVALID

		case http.StatusMisdirectedRequest:
			return EINVALID

		case http.StatusUnprocessableEntity:
			return EINVALID

		case http.StatusLocked:
			return EFORBIDDEN

		case http.StatusFailedDependency:
			return EINVALID

		case http.StatusTooEarly:
			return EINVALID

		case http.StatusUpgradeRequired:
			return EINVALID

		case http.StatusPreconditionRequired:
			return EINVALID

		case http.StatusTooManyRequests:
			return ETIMEOUT

		case http.StatusRequestHeaderFieldsTooLarge:
			return EINVALID

		case http.StatusUnavailableForLegalReasons:
			return EUNAVAILABLE

		case http.StatusBadGateway:
			return EUNAVAILABLE

		case http.StatusServiceUnavailable:
			return EUNAVAILABLE

		case http.StatusGatewayTimeout:
			return ETIMEOUT

		case http.StatusHTTPVersionNotSupported:
			return EINVALID
		}
	}
	if As(err, &gc) {
		switch gc.GRPCStatus().Code() {
		case codes.OK:
			return ""

		case codes.Canceled:
			return ECANCELED

		case codes.Unknown:
			return EUNAVAILABLE

		case codes.InvalidArgument:
			return EINVALID

		case codes.DeadlineExceeded:
			return ETIMEOUT

		case codes.NotFound:
			return ENOTFOUND

		case codes.AlreadyExists:
			return ECONFLICT

		case codes.PermissionDenied:
			return EFORBIDDEN

		case codes.ResourceExhausted:
			return EUNAVAILABLE

		case codes.FailedPrecondition:
			return EINVALID

		case codes.Aborted:
			return ECANCELED

		case codes.OutOfRange:
			return EINVALID

		case codes.Unimplemented:
			return EUNAVAILABLE

		case codes.Internal:
			return EINTERNAL

		case codes.Unavailable:
			return EUNAVAILABLE

		case codes.DataLoss:
			return ECONFLICT

		case codes.Unauthenticated:
			return EFORBIDDEN
		}
	}

	return EINTERNAL
}
