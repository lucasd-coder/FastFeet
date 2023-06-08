package monitor

import (
	"context"

	"github.com/lucasd-coder/user-manger-service/config"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"google.golang.org/grpc"
)

func RegisterOtel(ctx context.Context, cfg *config.Config) *sdktrace.TracerProvider {
	log := logger.FromContext(ctx)

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceName(cfg.Name),
		),
	)

	if err != nil {
		log.Errorf("fail creating OTLP trace resource: %v", err)
	}

	var traceExporter sdktrace.SpanExporter

	switch cfg.OpenTelemetry.Protocol {
	case "http":
		exp, err := registerExporterHTTP(ctx, cfg)
		if err != nil {
			log.Errorf("fail creating OTLP trace exporter: %w", err)
		}
		traceExporter = exp
	default:
		exp, err := registerExporterGRPC(ctx, cfg)
		if err != nil {
			log.Errorf("fail creating OTLP trace exporter: %w", err)
		}
		traceExporter = exp
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Info("creating OTLP trace exporter")

	return tracerProvider
}

func registerExporterGRPC(ctx context.Context, cfg *config.Config) (*otlptrace.Exporter, error) {
	conn := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OpenTelemetry.URL),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	export, err := otlptrace.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	return export, nil
}

func registerExporterHTTP(ctx context.Context, cfg *config.Config) (*otlptrace.Exporter, error) {
	return otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.OpenTelemetry.URL),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
}
