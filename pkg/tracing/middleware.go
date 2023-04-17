package tracing

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// OtelTracerMiddleware sends open telemetry traces to a configured exporter
func OtelTracerMiddleware(operation string, tp trace.TracerProvider, propagator propagation.TextMapPropagator) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return otelhttp.NewHandler(
			handler,
			operation,
			otelhttp.WithTracerProvider(tp),
			otelhttp.WithPropagators(propagator),
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
		)
	}
}

func spanNameFormatter(operation string, r *http.Request) string {
	return r.URL.Path
}
