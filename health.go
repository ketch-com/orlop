// Copyright (c) 2020 SwitchBit, Inc.
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

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

// HealthChecker provides the capability to check the health
type HealthChecker func(ctx context.Context, check string) (proto.Message, error)

// HealthHandler is a HTTP handler for checking health
type HealthHandler struct {
	checker HealthChecker
}

// ServeHTTP serves HTTP requests for `/healthz/`, optionally with a specific
// check appended
func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/healthz/") {
		check := strings.TrimPrefix(r.URL.Path, "/healthz/")
		if len(check) > 0 && h.checker != nil {
			var b []byte
			statusCode := 200

			s, err := h.checker(context.Background(), check)
			if err != nil {
				statusCode = 500

				b, _ = json.Marshal(&ErrorMessage{
					Code:    500,
					Error:   err.Error(),
					Message: err.Error(),
				})
			} else {
				b, _ = json.Marshal(s)
			}

			w.WriteHeader(statusCode)
			w.Write(b)
		}
	}
}
