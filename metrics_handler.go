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

func NewInstrumentedMetricHandlerFunc(next http.HandlerFunc) (*InstrumentedMetricHandler, error) {
	return NewInstrumentedMetricHandler(next)
}

func NewInstrumentedMetricHandler(next http.Handler) (*InstrumentedMetricHandler, error) {
	reg := prometheus.DefaultRegisterer

	inflightRequests := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "requests_in_flight",
		Help: "Current number of requests being served.",
	}, []string{"method", "route"})
	if err := reg.Register(inflightRequests); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			inflightRequests = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return nil, err
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
			return nil, err
		}
	}

	return &InstrumentedMetricHandler{
		requestDuration:  requestDuration,
		inflightRequests: inflightRequests,
		next:             next,
	}, nil
}

func (h InstrumentedMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.inflightRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
	defer h.inflightRequests.WithLabelValues(r.Method, r.URL.Path).Dec()

	m := httpsnoop.CaptureMetrics(h.next, w, r)
	h.requestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(m.Code)).Observe(m.Duration.Seconds())
}
