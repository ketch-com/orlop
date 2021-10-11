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

import (
	"net/http"
	"strings"
)

type HeaderOptions struct {
	AllowedOrigins []string
}

// DefaultHTTPHeaders is middleware to handle default HTTP headers
func DefaultHTTPHeaders(options HeaderOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isGRPCRequest := r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc")
			if isGRPCRequest {
				next.ServeHTTP(w, r)
				return
			}

			// Add CORS headers
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")

			if origin := r.Header.Get("Origin"); origin != "" {
				if len(options.AllowedOrigins) > 0 {
					for _, o := range options.AllowedOrigins {
						if o == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
						}
					}
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
					w.Header().Set("Access-Control-Allow-Methods", methods)
					return
				}
			}

			addSecurityHeaders(w, r)

			next.ServeHTTP(w, r)
		})
	}
}

func addSecurityHeaders(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		// Only on TLS per https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}

	addHeaderIfNotExists(w, "X-Frame-Options", "deny")
	addHeaderIfNotExists(w, "X-Content-Type-Options", "nosniff")
	addHeaderIfNotExists(w, "Content-Security-Policy", "default-src 'self'")
	addHeaderIfNotExists(w, "X-XSS-Protection", "1; mode=block")
}

func addHeaderIfNotExists(w http.ResponseWriter, headerKey string, value string) {
	if len(w.Header().Get(headerKey)) == 0 {
		w.Header().Add(http.CanonicalHeaderKey(headerKey), value)
	}
}
