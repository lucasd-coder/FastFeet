package publish

import (
	"context"
	"time"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/router-service/internal/shared"
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/trace"
	"gocloud.dev/pubsub"
)

type Published struct {
	opt *shared.Options
}

func NewPublished(opt *shared.Options) *Published {
	return &Published{
		opt: opt,
	}
}

func (p *Published) Send(ctx context.Context, msg *shared.Message) error {
	log := logger.FromContext(ctx)

	traceName := "gocloud.dev/pubsub/Topic.Send"
	tracer := otel.GetTracerProvider().Tracer(traceName)
	octrace.DefaultTracer = opencensus.NewTracer(tracer)

	commonAttrs := []attribute.KeyValue{
		attribute.String("queueURL", p.opt.TopicURL),
	}
	ctx, span := tracer.Start(ctx, "Topic.Send",
		trace.WithAttributes(commonAttrs...),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	client, err := NewClient(ctx, p.opt.TopicURL)
	if err != nil {
		span.RecordError(err)
		log.Errorf("error creating Publish client: %v", err)
	}

	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			span.RecordError(err)
			log.Fatalf("error client shutdown: %v", err)
		}
	}()

	m := &pubsub.Message{
		Body:     msg.Body,
		Metadata: msg.Metadata,
	}

	var er error
	for i := 0; i < p.opt.MaxRetries; i++ {
		er = client.Send(ctx, m)
		if er == nil {
			break
		}
		log.Errorf("error when trying to publish to queue with err: %v", er)

		if i == p.opt.MaxRetries-1 {
			log.Errorf("max retries exceeded, not publishing message anymore: %v", er)
			break
		}
		backOffTime := time.Duration(1+i) * p.opt.WaitingTime
		log.Infof("waiting %v before retrying", backOffTime)
		span.RecordError(err)
		time.Sleep(backOffTime)
	}
	return er
}
