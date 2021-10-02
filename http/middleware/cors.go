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

package middleware

import (
	"go.uber.org/fx"
	"net/http"
	"strings"
)

var (
	headers = strings.Join([]string{"Content-Type", "Accept", "Authorization"}, ",")
	methods = strings.Join([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}, ",")
)

type CORSParams struct {
	fx.In
}

// CORS is middleware to handle CORS headers
func CORS(params CORSParams) Middleware {
	return func(next http.Handler) http.Handler {
		return &CORSMiddleware{next}
	}
}

type CORSMiddleware struct {
	next http.Handler
}

func (h CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skip CORS for GRPC requests
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		h.next.ServeHTTP(w, r)
		return
	}

	if r.TLS != nil {
		// Only on TLS per https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}

	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Allow-Methods", methods)
			return
		}
	}

	h.next.ServeHTTP(w, r)
}
