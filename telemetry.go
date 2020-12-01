package orlop

import (
	"go.ketch.com/lib/orlop/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.GetTracerProvider().Tracer(version.Name, trace.WithInstrumentationVersion(version.Version))

var metrics = otel.GetMeterProvider().Meter(version.Name, metric.WithInstrumentationVersion(version.Version))
