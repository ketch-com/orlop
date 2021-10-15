package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/prometheus"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewPrometheusConfig() prometheus.Config {
	return prometheus.Config{
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
}

func NewPrometheusController(ctx context.Context, cfg prometheus.Config) (*controller.Controller, error) {
	res := resource.Environment()
	attributes := res.Attributes()

	res, err := resource.New(ctx,
		resource.WithSchemaURL(res.SchemaURL()),
		resource.WithAttributes(attributes...))
	if err != nil {
		return nil, err
	}

	return controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(cfg.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(res),
	), nil
}
