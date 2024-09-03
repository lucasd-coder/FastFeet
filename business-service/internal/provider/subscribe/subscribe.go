package subscribe

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/utils"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gocloud.dev/pubsub"
)

type Subscription struct {
	mux     sync.RWMutex
	handler func(ctx context.Context, m []byte) error
	opt     *queueoptions.Options
	metr    monitor.Metrics
	client  *pubsub.Subscription
}

func New(
	ctx context.Context,
	handler func(ctx context.Context, m []byte) error,
	opt *queueoptions.Options,
	metr monitor.Metrics) (*Subscription, error) {
	client, err := NewClient(ctx, opt)
	if err != nil {
		logger.FromContext(ctx).Errorf("error creating subscription client: %v, for queueURL",
			err, opt.QueueURL)
		return nil, err
	}

	return &Subscription{
		handler: handler,
		opt:     opt,
		metr:    metr,
		client:  client,
	}, nil
}

func (s *Subscription) Start(ctx context.Context) {
	s.mux.RLock()
	tracer := s.initializeTracer()
	s.mux.RUnlock()
	logDefault := logger.FromContext(ctx)

	logDefault.Infof("Subscription has been started.... for queueURL: %s", s.opt.QueueURL)

	commonAttrs := []attribute.KeyValue{
		attribute.String("queueURL", s.opt.QueueURL),
	}

	ctx, span := tracer.Start(ctx, "Subscription.Receive",
		trace.WithAttributes(commonAttrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	defer span.End()

	defer func() {
		if err := s.client.Shutdown(ctx); err != nil {
			span.RecordError(err)
			log.Fatalf("error client for queueURL: %s, shutdown: %v", s.opt.QueueURL, err)
		}
	}()

	msgChan := make(chan *pubsub.Message)

	go s.startReceivers(ctx, msgChan)

	var wg sync.WaitGroup
	wg.Add(s.opt.MaxConcurrentMessages)
	for i := 0; i < s.opt.MaxConcurrentMessages; i++ {
		go s.startProcess(ctx, &wg, msgChan)
	}
	wg.Wait()
	close(msgChan)
}

func (s *Subscription) startProcess(ctx context.Context, wg *sync.WaitGroup, msgChan chan *pubsub.Message) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logger.FromContext(ctx).Infof("context cancelled, stopping Subscription... for queueURL: %s", s.opt.QueueURL)
			runtime.Goexit()
			return
		case msg := <-msgChan:
			if err := s.processMessage(ctx, msg.Body); err != nil {
				logger.FromContext(ctx).Errorf("error processing message for queueURL: %s, err: %v", s.opt.QueueURL, err)
				if msg.Nackable() {
					defer msg.Nack()
				}
				return
			} else {
				msg.Ack()
			}
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
		err = s.handler(ctx, messages)
		if err == nil {
			span.SetStatus(codes.Ok, "Successfully Processing Message")
			s.createMetrics(monitor.OK, name, start)
			break
		}
		log.Errorf("error while handling message: %v", err)
		span.SetStatus(codes.Error, "Error Processing Message")
		span.RecordError(err)

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
		span.End()
	}
	return err
}

func (s *Subscription) startReceivers(ctx context.Context, m chan *pubsub.Message) {
	for i := 0; i < s.opt.NumberOfMessageReceivers; i++ {
		go s.receiveMessage(ctx, m)
	}
}

func (s *Subscription) receiveMessage(ctx context.Context, m chan *pubsub.Message) {
	if s.client == nil {
		return
	}

	log := logger.FromContext(ctx)
	start := time.Now()
	name := fmt.Sprintf("%s_receive", utils.ExtractQueueName(s.opt.QueueURL))

	log.Infof("start receive mensagens for queueURL: %s", s.opt.QueueURL)

	span := trace.SpanFromContext(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping receive... for queueURL %s", s.opt.QueueURL)
			return
		default:
			childCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			msg, err := s.client.Receive(childCtx)
			if err != nil {
				s.handleReceiveError(ctx, name, start, err)
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

	return tracer
}

func (s *Subscription) handleReceiveError(ctx context.Context, name string, start time.Time, err error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)

	span.RecordError(err)
	s.createMetrics(monitor.ERROR, name, start)
	log.Errorf("error receiving message for queueURL: %s, err: %v", s.opt.QueueURL, err)
	time.Sleep(s.opt.MaxReceiveMessage)

	client, err := NewClient(ctx, s.opt)
	if err != nil {
		span.RecordError(err)
		log.Errorf("error creating subscription client: %v, for queueURL", err, s.opt.QueueURL)
		return
	}
	s.updateClient(client)
}

func (s *Subscription) updateClient(client *pubsub.Subscription) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.client = client
}
