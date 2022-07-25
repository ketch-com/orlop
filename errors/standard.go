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

import "go.ketch.com/lib/orlop/v2/errors/internal"

// Application error codes.
const (
	ECONFLICT      = internal.ECONFLICT      // action cannot be performed
	EINTERNAL      = internal.EINTERNAL      // internal error
	EUNAVAILABLE   = internal.EUNAVAILABLE   // service is unavailable
	EINVALID       = internal.EINVALID       // validation failed
	ENOTFOUND      = internal.ENOTFOUND      // entity does not exist
	ETIMEOUT       = internal.ETIMEOUT       // operation timed out
	ECANCELED      = internal.ECANCELED      // operation canceled
	EFORBIDDEN     = internal.EFORBIDDEN     // operation is not authorized
	ECONFIGURATION = internal.ECONFIGURATION // configuration error
	EUNIMPLEMENTED = internal.EUNIMPLEMENTED // unimplemented error
)
