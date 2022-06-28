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
)

type messenger struct {
	error
	msg string
}

func (msgr messenger) Unwrap() error {
	return msgr.error
}

func (msgr messenger) UserMessage() string {
	return msgr.msg
}

// WithUserMessage adds a UserMessenger to err's error chain.
// If a status code has not previously been set,
// a default status of Bad Request (400) is added.
// Unlike pkg/errors, WithUserMessage will wrap nil error.
func WithUserMessage(err error, msg string) error {
	if err == nil {
		err = New("UserMessage<" + msg + ">")
	}
	return messenger{err, msg}
}

// WithUserMessagef calls fmt.Sprintf before calling WithUserMessage.
func WithUserMessagef(err error, format string, v ...any) error {
	return WithUserMessage(err, fmt.Sprintf(format, v...))
}

// UserMessage returns the user message associated with an error.
// If no message is found, it checks StatusCode and returns that message.
// Because the default status is 500, the default message is
// "Internal Server Error".
// If err is nil, it returns "".
func UserMessage(err error) string {
	var um internal.UserMessage

	if err == nil {
		return ""
	}
	if As(err, &um) {
		return um.UserMessage()
	}

	return http.StatusText(StatusCode(err))
}
