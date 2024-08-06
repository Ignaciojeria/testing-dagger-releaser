package observability

import (
	"testing-releaser/app/shared/configuration"
	"testing-releaser/app/shared/infrastructure/shutdown"
	"context"
	"fmt"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var Tracer = otel.Tracer("observability")

func init() {
	ioc.Registry(
		newTracerProvider,
		configuration.NewConf)
}

func newTracerProvider(conf configuration.Conf) error {
	ctx, cancel := context.WithCancel(context.Background())

	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		cancel()
		return fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.PROJECT_NAME),
			semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
		)),
	)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	shutdown.Handler(ctx, tp.Shutdown, time.Second*2, cancel)
	return nil
}
