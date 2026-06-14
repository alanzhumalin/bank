package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Init(ctx context.Context) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure())

	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("bank-api"))

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(res))

	otel.SetTracerProvider(tp)

	return tp, nil
}
