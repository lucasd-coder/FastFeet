package publish

import (
	"context"
	"time"

	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/pkg/logger"
	"gocloud.dev/pubsub"
)

type Published struct {
	cfg *config.Config
}

func NewPublished(cfg *config.Config) *Published {
	return &Published{
		cfg: cfg,
	}
}

func (p *Published) Send(ctx context.Context, msg *shared.Message) error {
	log := logger.FromContext(ctx)

	client, err := NewClient(ctx, p.cfg)
	if err != nil {
		log.Errorf("error creating Publish client: %v", err)
	}

	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Fatalf("error client shutdown: %v", err)
		}
	}()

	m := &pubsub.Message{
		Body:     msg.Body,
		Metadata: msg.Metadata,
	}

	var er error
	for i := 0; i < p.cfg.MaxRetries; i++ {
		er = client.Send(ctx, m)
		if er == nil {
			break
		}
		log.Errorf("error when trying to publish to queue with err: %v", er)

		if i == p.cfg.MaxRetries-1 {
			log.Errorf("max retries exceeded, not publishing message anymore: %v", er)
			break
		}
		backOffTime := time.Duration(1+i) * p.cfg.WaitingTime
		log.Infof("waiting %v before retrying", backOffTime)
		time.Sleep(backOffTime)
	}
	return er
}
