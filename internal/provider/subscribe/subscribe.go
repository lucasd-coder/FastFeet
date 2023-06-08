package subscribe

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lucasd-coder/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gocloud.dev/pubsub"
)

type Subscription struct {
	queueURL string
	handler  func(ctx context.Context, m []byte) error
	opt      *queueoptions.Options
}

func New(queueURL string, handler func(ctx context.Context, m []byte) error, opt *queueoptions.Options) *Subscription {
	return &Subscription{
		queueURL,
		handler,
		opt,
	}
}

func (s *Subscription) Start(ctx context.Context) {
	log := logger.FromContext(ctx)

	log.Infof("Subscription has been started.... for queueURL: %s", s.queueURL)

	client, err := NewClient(ctx, s.queueURL)
	if err != nil {
		log.Errorf("error creating subscription client: %v, for queueURL", err, s.queueURL)
	}

	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Fatalf("error client for queueURL: %s, shutdown: %v", s.queueURL, err)
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
			log.Infof("context cancelled, stopping Subscription... for queueURL: %s", s.queueURL)
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
					log.Errorf("error processing message for queueURL: %s, err: %v", s.queueURL, err)
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
	log.Infof("start process mensagens for queueURL: %s", s.queueURL)

	spanName := fmt.Sprintf("Processing-%s", s.extractQueueName(s.queueURL))
	traceName := "ProcessingMessage"

	tracer := otel.GetTracerProvider().Tracer(traceName)

	commonAttrs := []attribute.KeyValue{
		attribute.String("queueURL", s.queueURL),
	}

	newCtx, span := tracer.Start(ctx, spanName,
		trace.WithAttributes(commonAttrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			span.SetStatus(codes.Error, "recovered from panic")
			log.Errorf("recovered from panic: %v", r)
		}
	}()

	var err error
	for i := 0; i < s.opt.MaxRetries; i++ {
		_, iSpan := tracer.Start(newCtx, fmt.Sprintf(" ProcessingMessage MaxRetries-%d", i))
		err = s.handler(ctx, messages)
		if err == nil {
			iSpan.SetStatus(codes.Ok, "Successfully Processing Message")
			break
		}
		log.Errorf("error while handling message: %v", err)
		iSpan.SetStatus(codes.Error, "Error Processing Message")
		iSpan.RecordError(err)

		if i == s.opt.MaxRetries-1 {
			log.Errorf("max retries exceeded, not processing message anymore: %v", err)
			err = nil
			break
		}

		backOffTime := time.Duration(1+i) * s.opt.WaitingTime
		log.Infof("waiting %v before retrying", backOffTime)
		iSpan.End()
		time.Sleep(backOffTime)
	}
	return err
}

func (s *Subscription) startReceivers(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	for i := 0; i < s.opt.NumberOfMessageReceivers; i++ {
		go s.receiveMessage(ctx, client, m)
	}
}

func (s *Subscription) receiveMessage(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	log := logger.FromContext(ctx)
	log.Infof("start receive mensagens for queueURL: %s", s.queueURL)

	retry := s.opt.MaxReceiveMessage
	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping receive... for queueURL: %s", s.queueURL)
			return
		default:
			childCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			msg, err := client.Receive(childCtx)
			if err != nil {
				log.Errorf("error receiving message for queueURL: %s, err: %v", s.queueURL, err)
				time.Sleep(retry)
				continue
			}

			if len(msg.Body) > 0 {
				m <- msg
			}

			s.applyBackPressure()
		}
	}
}

func (s *Subscription) applyBackPressure() {
	time.Sleep(s.opt.PollDelay)
}

func (s *Subscription) extractQueueName(queueURL string) string {
	index := 2
	parts := strings.SplitN(queueURL, "rabbit://", index)

	queueName := parts[1]

	return queueName
}
