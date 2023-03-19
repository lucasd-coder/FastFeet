package subscribe

import (
	"context"

	// revive
	_ "gocloud.dev/pubsub/rabbitpubsub"

	"github.com/lucasd-coder/business-service/config"
	"gocloud.dev/pubsub"
)

func NewClient(ctx context.Context, cfg *config.Config) (*pubsub.Subscription, error) {
	queueURL := cfg.Integration.RabbitMQ.URL
	client, err := pubsub.OpenSubscription(ctx, queueURL)
	return client, err
}
