package tracing

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

const (
	OtelDefaultEndpoint = "otel-collector:4317"
	ConsoleEndpoint     = "stdout"
)

// InitTracerProvider creates a new [trace.TracerProvider] with default Otel gRPC [otlptrace.Client]
// and telemetry [resource.Resource].
// The TraceProvider is set as the global TraceProvider of the Otel SDK.
// It is also returned for a later shutdown.
func InitTracerProvider(ctx context.Context, serviceName, endpoint string) (*trace.TracerProvider, error) {
	exp, err := newSpanExporterOnEndpoint(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	r, err := newResource(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	tp := NewTracerProvider(exp, r)

	// set global propagator to trace context (the default is no-op).
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		))
	otel.SetTracerProvider(tp)

	return tp, nil
}

// NewTraceProvider registers a [otlptrace.Exporter] and a telemetry [resource.Resource]
// with a [trace.TracerProvider].
// The Exporter will receive spans in batches to increase performance.
func NewTracerProvider(exp trace.SpanExporter, r *resource.Resource) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithSampler(newSamplerOnEnv()),
		trace.WithResource(r),
		trace.WithBatcher(
			exp,
			trace.WithMaxQueueSize(3000)))
}

// ShutdownTracerProvider flushes the span processors' content to the exporter and
// shuts down the span processors registered with the TracerProvider; resources will be freed up.
func FlushAndShutdownTracerProvider(ctx context.Context, tp *trace.TracerProvider) error {
	return tp.Shutdown(ctx)
}

// newSpanExporterOnEndpoint creates a new Exporter,depending on the endpoint,
// where traces are sent to.
// Possible endpoints are currently stdout and the OpenTelemetry collector.
// If the otel collector is chosen, gRPC will be used as transport as of now.
// Returned is the [trace.SpanExporter] interface type.
//
// Todo: Add proper constants for possible exporters and transports
func newSpanExporterOnEndpoint(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	endpoint = checkEndpoint(endpoint)

	switch endpoint {
	case ConsoleEndpoint:
		return newConsoleExporter()
	default:
		return newgRPCTraceExporter(ctx, endpoint)
	}
}

func newgRPCTraceExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	return otlptrace.New(ctx, newgRPCClient(endpoint))
}

// ?Can we also get logging?
// Todo: To do proper timeout handling if the server isn't reachable,
// we must use our custom [grpc.ClientConn] created with [grpc.DialContext]
// together with [context.ContextWithTImeout].
// [grpc.WithTImeout] is deprecated anyway.
func newgRPCClient(endpoint string) otlptrace.Client {
	endpoint = checkEndpoint(endpoint)
	return otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock(), grpc.WithTimeout(5*time.Second)),
	)
}

func newConsoleExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithoutTimestamps())
}

// Todo: Have env var or so that returns fractional sampler
func newSamplerOnEnv() trace.Sampler {
	// return trace.AlwaysSample()

	// More options are possible
	return trace.ParentBased(trace.TraceIDRatioBased(0.7))
}

func newResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
}

func checkEndpoint(ep string) string {
	if ep == "" {
		var ok bool
		ep, ok = os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if !ok {
			ep = OtelDefaultEndpoint
		}
	}
	return ep
}
