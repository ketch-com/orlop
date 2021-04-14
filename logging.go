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
	"github.com/felixge/httpsnoop"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/log"
	"net/http"
)

type loggingMiddleware struct {
	cfg  HttpLoggingConfig
	next http.Handler
}

func (l *loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := httpsnoop.CaptureMetrics(l.next, w, r)

	headers := make(map[string][]string)
	for _, header := range l.cfg.Headers {
		headers[header] = r.Header[header]
	}

	fields := logrus.Fields{
		"status":        m.Code,
		"duration":      m.Duration.String(),
		"bytes":         m.Written,
		"method":        r.Method,
		"proto":         r.Proto,
		"contentLength": r.ContentLength,
		"host":          r.Host,
		"remoteAddr":    r.RemoteAddr,
		"userAgent":     r.UserAgent(),
		"headers":       headers,
	}

	log.WithFields(fields).Info(r.URL.Path)
}

// Logging is middleware to log each HTTP request
func Logging(cfg HttpLoggingConfig) func(http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return &loggingMiddleware{
			cfg:  cfg,
			next: next,
		}
	}
}
