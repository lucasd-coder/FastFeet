package monitor

import (
	"context"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

type Option struct {
	ServiceName string
	Protocol    string
	URL         string
}

func RegisterOtel(ctx context.Context, opt *Option) (*sdktrace.TracerProvider, error) {
	log := logger.FromContext(ctx)

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceName(opt.ServiceName),
		),
	)

	if err != nil {
		log.Errorf("fail creating OTLP trace resource: %v", err)
		return nil, err
	}

	var traceExporter sdktrace.SpanExporter

	switch opt.Protocol {
	case "http":
		exp, err := registerExporterHTTP(ctx, opt)
		if err != nil {
			log.Errorf("fail creating OTLP trace exporter: %w", err)
			return nil, err
		}
		traceExporter = exp
	default:
		exp, err := registerExporterGRPC(ctx, opt)
		if err != nil {
			log.Errorf("fail creating OTLP trace exporter: %w", err)
			return nil, err
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

	return tracerProvider, nil
}

func registerExporterGRPC(ctx context.Context, opt *Option) (*otlptrace.Exporter, error) {
	conn := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(opt.URL))

	export, err := otlptrace.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	return export, nil
}

func registerExporterHTTP(ctx context.Context, opt *Option) (*otlptrace.Exporter, error) {
	return otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(opt.URL),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
}
