package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

type OTel struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func New(ctx context.Context, serviceName string) (*OTel, error) {
	var tp *sdktrace.TracerProvider
	var err error
	tp, err = createOtelProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(serviceName)

	return &OTel{provider: tp, tracer: tracer}, nil
}

func createOtelProvider(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String("production"),
		)),
	)
	return tp, nil
}

func (o *OTel) Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, Span) {
	if len(opts) == 0 {
		return o.tracer.Start(ctx, name)
	}
	return o.tracer.Start(ctx, name, opts[0])
}

func (o *OTel) Shutdown(ctx context.Context) {
	o.provider.Shutdown(ctx)
}
