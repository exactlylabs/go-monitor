package gcpexporters

import (
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/exactlylabs/go-errors/pkg/errors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewTraceExporter returns an exporter that sends traces to Google Cloud Trace from the configured credential's project
func NewTraceExporter() (sdktrace.SpanExporter, error) {
	exporter, err := texporter.New()
	if err != nil {
		return nil, errors.Wrap(err, "New")
	}
	return exporter, nil
}
