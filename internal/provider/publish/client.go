package publish

import (
	"context"

	// revive
	_ "gocloud.dev/pubsub/rabbitpubsub"

	"github.com/lucasd-coder/router-service/config"
	"gocloud.dev/pubsub"
)

func NewClient(ctx context.Context, cfg *config.Config) (*pubsub.Topic, error) {
	queueURL := cfg.Integration.RabbitMQ.URL

	client, err := pubsub.OpenTopic(ctx, queueURL)
	return client, err
}
