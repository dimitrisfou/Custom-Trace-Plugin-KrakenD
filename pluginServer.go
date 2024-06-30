package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

// Initialization function to indicate the plugin is loaded
func init() {
	fmt.Println("TraceSpanGenerator plugin is loaded!")
}

func main() {}

// HandlerRegisterer is a global variable to register handlers
var HandlerRegisterer registrable = registrable("TraceSpanGenerator")

// registrable is a custom type used for handler registration
type registrable string

const (
	pluginName = "TraceSpanGenerator"
	nameTracer = "Gateway-API"
)

// RegisterHandlers registers the HTTP handlers with tracing
func (req registrable) RegisterHandlers(f func(
	name string,
	handler func(
		context.Context,
		map[string]interface{},
		http.Handler) (http.Handler, error),
)) {
	f(pluginName, req.registerHandlers)
}

// newTraceProvider initializes and returns a new trace provider
func newTraceProvider() (*trace.TracerProvider, error) {
	traceProvider := trace.NewTracerProvider()
	return traceProvider, nil
}

// Global tracer instance
var tracer = otel.Tracer(nameTracer)

// registerHandlers wraps the HTTP handler with tracing functionality
func (req registrable) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {

	// Get the TextMapPropagator instance from the OpenTelemetry SDK
	propagator := propagation.TraceContext{}

	// Set the global propagator
	otel.SetTextMapPropagator(propagator)

	// Set up trace provider
	tracerProvider, err := newTraceProvider()
	if err != nil {
		fmt.Println("Error setting up trace provider:", err)
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Start a new span for the incoming request with root kind
		newCtx, span := tracer.Start(
			ctx,
			"Root-Server-Span",
			otelTrace.WithSpanKind(otelTrace.SpanKindServer),
		)
		defer span.End() // Ensure the span is ended

		// Log the current trace and span IDs
		spanContext := span.SpanContext()
		fmt.Printf("TraceID: %s, SpanID: %s\n", spanContext.TraceID(), spanContext.SpanID())

		// Inject the span context into the HTTP headers
		carrier := propagation.HeaderCarrier(req.Header)
		otel.GetTextMapPropagator().Inject(newCtx, carrier)

		// Call the next handler in the chain with the new context
		handler.ServeHTTP(w, req.WithContext(newCtx))
	}), nil
}
