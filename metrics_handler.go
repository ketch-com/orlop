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
	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type InstrumentedMetricHandler struct {
	requestDuration  *prometheus.HistogramVec
	inflightRequests *prometheus.GaugeVec
	next             http.Handler
}

func Metrics(next http.Handler) http.Handler {
	reg := prometheus.DefaultRegisterer

	inflightRequests := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "requests_in_flight",
		Help: "Current number of requests being served.",
	}, []string{"method", "route"})
	if err := reg.Register(inflightRequests); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			inflightRequests = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return nil
		}
	}

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})
	if err := reg.Register(requestDuration); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			requestDuration = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			return nil
		}
	}

	return &InstrumentedMetricHandler{
		requestDuration:  requestDuration,
		inflightRequests: inflightRequests,
		next:             next,
	}
}

func (h InstrumentedMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.inflightRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
	defer h.inflightRequests.WithLabelValues(r.Method, r.URL.Path).Dec()

	m := httpsnoop.CaptureMetrics(h.next, w, r)
	h.requestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(m.Code)).Observe(m.Duration.Seconds())
}
