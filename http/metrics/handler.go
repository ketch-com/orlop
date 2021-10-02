package metrics

import (
	"go.ketch.com/lib/orlop/http/routing"
	"go.ketch.com/lib/orlop/logging"
	"go.opentelemetry.io/otel/exporters/prometheus"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Logger logging.Logger
}

// NewHandler creates a new MetricsHandler
func NewHandler(params Params) fx.Annotated {
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{
			0.005,
			0.01,
			0.025,
			0.05,
			0.1,
			0.25,
			0.5,
			1,
			10,
			2.5,
			5,
		},
	}

	c := controller.New(
		processor.New(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
	)

	exporter, err := prometheus.New(
		config,
		c,
	)
	if err != nil {
		params.Logger.Fatal(err)
	}

	return routing.GET("/metrics", exporter)
}
