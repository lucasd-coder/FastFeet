package subscribe

import (
	"context"

	// revive
	_ "gocloud.dev/pubsub/rabbitpubsub"

	"gocloud.dev/pubsub"
)

func NewClient(ctx context.Context, queueURL string) (*pubsub.Subscription, error) {
	client, err := pubsub.OpenSubscription(ctx, queueURL)
	return client, err
}
