package orlop

import (
	"go.ketch.com/lib/orlop/version"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

var tracer = global.TracerProvider().Tracer(version.Name, trace.WithInstrumentationVersion(version.Version))

var metrics = global.MeterProvider().Meter(version.Name, metric.WithInstrumentationVersion(version.Version))
