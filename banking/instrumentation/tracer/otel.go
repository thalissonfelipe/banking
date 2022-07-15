package tracer

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func New() (closer func() error, err error) {
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(os.Stderr),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating stdout exporter: %w", err)
	}

	tracer := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("banking"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			semconv.ServiceInstanceIDKey.String(uuid.NewString()),
		)),
		trace.WithSampler(trace.AlwaysSample()),
	)

	closer = func() error {
		if err := exporter.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("shutting down exporter: %w", err)
		}

		if err := tracer.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("shuttding down tracer: %w", err)
		}

		return nil
	}

	otel.SetTracerProvider(tracer)
	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))

	return closer, nil
}
