package subscribe

import (
	"context"

	"github.com/lucasd-coder/business-service/internal/shared/queueoptions"
	// revive
	_ "gocloud.dev/pubsub/rabbitpubsub"

	"gocloud.dev/pubsub"
)

func NewClient(ctx context.Context, opt *queueoptions.Options) (*pubsub.Subscription, error) {
	client, err := pubsub.OpenSubscription(ctx, opt.QueueURL)
	return client, err
}
