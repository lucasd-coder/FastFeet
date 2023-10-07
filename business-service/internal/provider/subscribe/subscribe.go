package subscribe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/utils"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"

	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gocloud.dev/pubsub"
)

var (
	mutex sync.Mutex
)

type Subscription struct {
	handler func(ctx context.Context, m []byte) error
	opt     *queueoptions.Options
	metr    monitor.Metrics
}

func New(
	handler func(ctx context.Context, m []byte) error,
	opt *queueoptions.Options,
	metr monitor.Metrics) *Subscription {
	return &Subscription{
		handler,
		opt,
		metr,
	}
}

func (s *Subscription) Start(ctx context.Context) {
	mutex.Lock()
	tracer := s.initializeTracer()
	mutex.Unlock()
	log := logger.FromContext(ctx)

	log.Infof("Subscription has been started.... for queueURL: %s", s.opt.QueueURL)

	commonAttrs := []attribute.KeyValue{
		attribute.String("queueURL", s.opt.QueueURL),
	}

	ctx, span := tracer.Start(ctx, "Subscription.Receive",
		trace.WithAttributes(commonAttrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	defer span.End()

	client, err := NewClient(ctx, s.opt)
	if err != nil {
		span.RecordError(err)
		log.Errorf("error creating subscription client: %v, for queueURL", err, s.opt.QueueURL)
	}

	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			span.RecordError(err)
			log.Fatalf("error client for queueURL: %s, shutdown: %v", s.opt.QueueURL, err)
		}
	}()

	var wg sync.WaitGroup
	sem := make(chan struct{}, s.opt.MaxConcurrentMessages)

	s.start(ctx, client, &wg, sem)
	wg.Wait()
	close(sem)
}

func (s *Subscription) start(ctx context.Context, client *pubsub.Subscription, wg *sync.WaitGroup, sem chan struct{}) {
	log := logger.FromContext(ctx)

	msgChan := make(chan *pubsub.Message)

	s.startReceivers(ctx, client, msgChan)

	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping Subscription... for queueURL: %s", s.opt.QueueURL)
			return
		case msg := <-msgChan:
			sem <- struct{}{}
			wg.Add(1)
			go func(ctx context.Context, currentMsg *pubsub.Message) {
				defer func() {
					<-sem
					wg.Done()
				}()

				if err := s.processMessage(ctx, currentMsg.Body); err != nil {
					log.Errorf("error processing message for queueURL: %s, err: %v", s.opt.QueueURL, err)
					if currentMsg.Nackable() {
						defer currentMsg.Nack()
					}
					return
				}
				defer currentMsg.Ack()
			}(ctx, msg)
		}
	}
}

func (s *Subscription) processMessage(ctx context.Context, messages []byte) error {
	log := logger.FromContext(ctx)
	start := time.Now()
	name := fmt.Sprintf("%s_consumed", utils.ExtractQueueName(s.opt.QueueURL))
	log.Infof("start process mensagens for queueURL: %s", s.opt.QueueURL)

	spanName := fmt.Sprintf("Processing-%s", utils.ExtractQueueName(s.opt.QueueURL))

	traceName := "Processing-Message"

	tracer := otel.GetTracerProvider().Tracer(traceName)

	commonAttrs := []attribute.KeyValue{
		attribute.String("queueURL", s.opt.QueueURL),
	}

	ctx, span := tracer.Start(ctx, spanName,
		trace.WithAttributes(commonAttrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			span.SetStatus(codes.Error, "recovered from panic")
			s.createMetrics(monitor.ERROR, name, start)
			log.Errorf("recovered from panic: %v", r)
		}
	}()

	var err error
	for i := 0; i < s.opt.MaxRetries; i++ {
		ctx, iSpan := tracer.Start(ctx, fmt.Sprintf(" ProcessingMessage MaxRetries-%d", i))
		err = s.handler(ctx, messages)
		if err == nil {
			iSpan.SetStatus(codes.Ok, "Successfully Processing Message")
			s.createMetrics(monitor.OK, name, start)
			break
		}
		log.Errorf("error while handling message: %v", err)
		iSpan.SetStatus(codes.Error, "Error Processing Message")
		iSpan.RecordError(err)

		if i == s.opt.MaxRetries-1 {
			log.Errorf("max retries exceeded, not processing message anymore: %v", err)
			s.createMetrics(monitor.ERROR, name, start)
			err = nil
			break
		}
		s.createMetrics(monitor.ERROR, name, start)
		backOffTime := time.Duration(1+i) * s.opt.WaitingTime
		log.Infof("waiting %v before retrying", backOffTime)
		time.Sleep(backOffTime)
		iSpan.End()
	}
	return err
}

func (s *Subscription) startReceivers(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	mutex.Lock()
	defer mutex.Unlock()
	for i := 0; i < s.opt.NumberOfMessageReceivers; i++ {
		go s.receiveMessage(ctx, client, m)
	}
}

func (s *Subscription) receiveMessage(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	if client != nil {
		log := logger.FromContext(ctx)
		start := time.Now()
		name := fmt.Sprintf("%s_receive", utils.ExtractQueueName(s.opt.QueueURL))

		log.Infof("start receive mensagens for queueURL: %s", s.opt.QueueURL)

		retry := s.opt.MaxReceiveMessage
		span := trace.SpanFromContext(ctx)
		for {
			select {
			case <-ctx.Done():
				log.Infof("context cancelled, stopping receive... for queueURL %s", s.opt.QueueURL)
				return
			default:
				childCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				msg, err := client.Receive(childCtx)
				if err != nil {
					span.RecordError(err)
					s.createMetrics(monitor.ERROR, name, start)
					log.Errorf("error receiving message for queueURL: %s, err: %v", s.opt.QueueURL, err)
					time.Sleep(retry)
					client, err = NewClient(ctx, s.opt)
					if err != nil {
						span.RecordError(err)
						log.Errorf("error creating subscription client: %v, for queueURL", err, s.opt.QueueURL)
					}
					continue
				}

				if msg != nil && len(msg.Body) > 0 {
					s.createMetrics(monitor.OK, name, start)
					m <- msg
				}
				s.applyBackPressure()
				span.End()
			}
		}
	}
}

func (s *Subscription) applyBackPressure() {
	time.Sleep(s.opt.PollDelay)
}
func (s *Subscription) createMetrics(status string, queueName string, observeTime time.Time) {
	s.metr.ObserveResponseTime(status, queueName, time.Since(observeTime).Seconds())
	s.metr.IncHits(monitor.OK, queueName)
}

func (s *Subscription) initializeTracer() trace.Tracer {
	traceName := "gocloud.dev/pubsub/Subscription.Receive"
	tracer := otel.GetTracerProvider().Tracer(traceName)
	octrace.DefaultTracer = opencensus.NewTracer(tracer)

	return tracer
}
