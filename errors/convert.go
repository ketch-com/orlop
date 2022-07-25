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
	"io/fs"
	"os"
)

// Converter is an interface for converting errors
type Converter func(err error) (error, bool)

// Convert converts an external error into a standard error
func Convert(err error) error {
	if err == nil {
		return nil
	}

	for _, converter := range internal.Converters {
		newErr, ok := converter(err)
		if ok {
			return newErr
		}
	}

	return WithCode(err, EINTERNAL)
}

// RegisterConverter registers a converter
func RegisterConverter(converter ...Converter) {
	for _, c := range converter {
		internal.RegisterConverter(c)
	}
}

func convertOSError(err error) (error, bool) {
	switch {
	case os.IsPermission(err):
		return Forbidden(err), true

	case os.IsExist(err):
		return Conflict(err), true

	case os.IsNotExist(err):
		return NotFound(err), true

	case err == fs.ErrInvalid:
		return Invalid(err), true

	case err == os.ErrDeadlineExceeded:
		return Timeout(err), true

	case err == os.ErrNoDeadline:
		return Configuration(err), true

	default:
		return nil, false
	}
}

func convertInterfaces(err error) (error, bool) {
	var ec internal.ErrorCode
	var sc internal.StatusCode
	var s internal.Status
	var gc internal.GRPCStatus
	var to internal.Timeout

	switch {
	case As(err, &ec):
		return ec, true

	case As(err, &sc):
		return WithStatusCode(err, sc.StatusCode()), true

	case As(err, &s):
		return WithStatusCode(err, s.Status()), true

	case As(err, &gc):
		return WithCode(err, internal.GrpcToStandard[gc.GRPCStatus().Code()]), true

	case As(err, &to):
		return Timeout(err), true

	default:
		return nil, false
	}
}

func init() {
	RegisterConverter(convertOSError, convertInterfaces)
}
