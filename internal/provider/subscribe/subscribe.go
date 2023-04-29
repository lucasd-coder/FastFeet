package subscribe

import (
	"context"
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

	var wg sync.WaitGroup

	s.start(ctx, client, &wg)

	wg.Wait()
}

func (s *Subscription) start(ctx context.Context, client *pubsub.Subscription, wg *sync.WaitGroup) {
	log := logger.FromContext(ctx)

	msgChan := make(chan *pubsub.Message)

	sem := make(chan struct{}, s.cfg.MaxConcurrentMessages)

	s.startReceivers(ctx, client, msgChan)

	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled, stopping Subscription...")
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
					log.Errorf("error processing message: %v", err)
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

		backOffTime := time.Duration(1+i) * s.cfg.WaitingTime
		log.Infof("waiting %v before retrying", backOffTime)
		time.Sleep(backOffTime)
	}
	return err
}

func (s *Subscription) startReceivers(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	for i := 0; i < s.cfg.NumberOfMessageReceivers; i++ {
		go s.receiveMessage(ctx, client, m)
	}
}

func (s *Subscription) receiveMessage(ctx context.Context, client *pubsub.Subscription, m chan *pubsub.Message) {
	log := logger.FromContext(ctx)
	log.Info("start receive mensagens")

	retry := s.cfg.MaxReceiveMessage
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

func (s *Subscription) applyBackPressure() {
	time.Sleep(s.cfg.PollDelay)
}
