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

package middleware

import (
	"github.com/felixge/httpsnoop"
	"go.ketch.com/lib/orlop/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"net/http"
)

type MetricsParams struct {
	fx.In

	Logger  logging.Logger
	Metrics metric.Meter
}

// Metrics is middleware for handling metrics
func Metrics(params MetricsParams) Middleware {
	return func(next http.Handler) http.Handler {
		inflightRequests, err := params.Metrics.NewInt64UpDownCounter("requests.in.flight", metric.WithUnit("requests"))
		if err != nil {
			params.Logger.Fatal(err)
		}

		requestDuration, err := params.Metrics.NewFloat64Histogram("request.duration.seconds", metric.WithUnit("s"))
		if err != nil {
			params.Logger.Fatal(err)
		}

		return &MetricsMiddleware{
			next:             next,
			inflightRequests: inflightRequests,
			requestDuration:  requestDuration,
		}
	}
}

type MetricsMiddleware struct {
	next             http.Handler
	inflightRequests metric.Int64UpDownCounter
	requestDuration  metric.Float64Histogram
}

func (m MetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := attribute.String("method", r.Method)
	route := attribute.String("route", r.URL.Path)

	m.inflightRequests.Add(r.Context(), 1, method, route)
	defer m.inflightRequests.Add(r.Context(), -1, method, route)

	metrics := httpsnoop.CaptureMetrics(m.next, w, r)
	m.requestDuration.Record(r.Context(), metrics.Duration.Seconds(), method, route, attribute.Int("status_code", metrics.Code))
}
