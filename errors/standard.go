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

type ErrorCode string

// Application error codes.
const (
	ECONFLICT    ErrorCode = "conflict"    // action cannot be performed - permanent
	EINTERNAL    ErrorCode = "internal"    // internal error - retryable (probably)
	EUNAVAILABLE ErrorCode = "unavailable" // service is unavailable - permanent
	EINVALID     ErrorCode = "invalid"     // validation failed - retryable (after fixing the input)
	ENOTFOUND    ErrorCode = "not_found"   // entity does not exist - permanent
	ETIMEOUT     ErrorCode = "timeout"     // operation timed out - retryable
	EFORBIDDEN   ErrorCode = "forbidden"   // operation is not authorized
)
