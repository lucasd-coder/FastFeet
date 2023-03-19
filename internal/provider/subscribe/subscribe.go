package subscribe

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"gocloud.dev/pubsub"
)

type Subscription struct {
	handler func(ctx context.Context, m []byte) error
	cfg     *config.Config
}

func New(handler func(ctx context.Context, m []byte) error, cfg *config.Config) *Subscription {
	return &Subscription{
		handler,
		cfg,
	}
}

func (s *Subscription) Start(ctx context.Context) {
	log := logger.FromContext(ctx)

	log.Info("Subscription has been started...")

	client, err := NewClient(ctx, s.cfg)
	if err != nil {
		log.Errorf("error creating Subscription client: %v", err)
	}

	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Fatalf("error client shutdown: %v", err)
		}
	}()

	msgChan := make(chan *pubsub.Message)

	sem := make(chan struct{}, s.cfg.MaxConcurrentMessages)

	var wg sync.WaitGroup

	go s.receive(ctx, client, msgChan)

	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping Subscription...")
			wg.Wait()
			return
		case msg := <-msgChan:
			currentMsg := msg
			sem <- struct{}{}
			wg.Add(1)
			go func(ctx context.Context) {
				defer func() {
					<-sem
					wg.Done()
				}()
				if err := s.process(ctx, currentMsg.Body); err != nil {
					log.Errorf("error processing message: %v", err)
					if currentMsg.Nackable() {
						defer currentMsg.Nack()
					}
					return
				}
				defer currentMsg.Ack()
			}(ctx)
		}
	}
}

func (s *Subscription) receive(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	log := logger.FromContext(ctx)
	log.Info("start receive mensagens")

	retry, err := time.ParseDuration(s.cfg.MaxReceiveMessage)
	if err != nil {
		log.Errorf("err parse duration to max receive message: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping receive...")
			return
		default:
			childCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			msg, err := client.Receive(childCtx)
			if err != nil {
				log.Errorf("error receiving message: %v", err)
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

func (s *Subscription) process(ctx context.Context, messages []byte) error {
	log := logger.FromContext(ctx)
	log.Info("start process mensagens")

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from panic: %v", r)
		}
	}()

	var err error
	for i := 0; i < s.cfg.MaxRetries; i++ {
		err = s.handler(ctx, messages)
		if err == nil {
			break
		}
		log.Errorf("error while handling message: %v", err)

		if i == s.cfg.MaxRetries-1 {
			log.Errorf("max retries exceeded, not processing message anymore: %v", err)
			err = nil
			break
		}

		backOffTime := time.Duration(math.Pow(s.cfg.WaitingTime, float64(i))) * time.Second
		log.Infof("waiting %v before retrying", backOffTime)
		time.Sleep(backOffTime)
	}
	return err
}

func (s *Subscription) applyBackPressure() {
	time.Sleep(time.Millisecond * time.Duration(s.cfg.PollDelayInMilliseconds))
}
