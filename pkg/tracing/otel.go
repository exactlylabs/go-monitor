package tracing

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// NewTracerProvider configures a Tracer that will only trace if the parent has tracing enabled
// or if based on defaultSampler, when the parent has no tracing configured
func NewTracerProvider(serviceName string, defaultSampler sdktrace.Sampler, exporter sdktrace.SpanExporter) (trace.TracerProvider, error) {
	svcId, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "tracing.NewTracerProvider NewV4")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(defaultSampler)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceInstanceIDKey.String((svcId.String())),
			),
		),
	)
	return tp, nil
}
