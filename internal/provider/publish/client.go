package publish

import (
	"context"

	// revive
	_ "gocloud.dev/pubsub/rabbitpubsub"

	"gocloud.dev/pubsub"
)

func NewClient(ctx context.Context, queueURL string) (*pubsub.Topic, error) {
	client, err := pubsub.OpenTopic(ctx, queueURL)
	return client, err
}
