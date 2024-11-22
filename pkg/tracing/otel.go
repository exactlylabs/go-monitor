package tracing

import (
	"github.com/exactlylabs/go-errors/pkg/errors"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTracerProvider configures a Tracer that will only trace if the parent has tracing enabled
// or if based on defaultSampler, when the parent has no tracing configured
func NewTracerProvider(serviceName string, defaultSampler trace.Sampler, exporter trace.SpanExporter) (*trace.TracerProvider, error) {
	svcId, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "NewV4")
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.ParentBased(defaultSampler)),
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceInstanceIDKey.String((svcId.String())),
			),
		),
	)
	return tp, nil
}
