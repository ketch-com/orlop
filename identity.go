// Copyright (c) 2020 Ketch Kloud, Inc.
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

package orlop

import "context"

// Identity returns the given input
func Identity[T any](in T) T {
	return in
}

// IdentityError returns the given input
func IdentityError[T any](in T) (T, error) {
	return in, nil
}

// IdentityContext returns the given input with a context as the first parameter
func IdentityContext[T any](_ context.Context, in T) T {
	return in
}

// IdentityContextError returns the given input with a context as the first parameter with a nil error
func IdentityContextError[T any](_ context.Context, in T) (T, error) {
	return in, nil
}

// IgnoreError returns the given input and ignores the error (use at your own peril
func IgnoreError[T any](in T, _ error) T {
	return in
}

// PanicOnError panics on error and returns the given input if not an error
func PanicOnError[T any](in T, err error) T {
	if err != nil {
		panic(err)
	}
	return in
}
