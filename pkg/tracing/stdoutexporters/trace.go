package stdoutexporters

import (
	"github.com/exactlylabs/go-errors/pkg/errors"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewTraceExporter() (sdktrace.SpanExporter, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, errors.Wrap(err, "New")
	}
	return exporter, nil
}
