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

package health

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Handler is an HTTP handler for checking health
type Handler struct {
	checks []Check
}

// ServeHTTP serves HTTP requests for `/healthz/`, optionally with a specific
// check appended
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	requestedCheck := chi.URLParam(r, "check")

	results := make(map[string]interface{})
	for _, check := range h.checks {
		if len(requestedCheck) > 0 && check.Name != requestedCheck {
			continue
		}

		out, err := check.Checker(r.Context())
		if err != nil {
			statusCode = http.StatusInternalServerError
			out = err.Error()
		} else if out == nil {
			out = "ok"
		}

		results[check.Name] = out
	}

	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(results)
}
