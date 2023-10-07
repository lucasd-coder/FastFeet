package cache

import (
	"context"
	"time"

	"github.com/lucasd-coder/business-service/internal/shared/codec"
	"github.com/redis/go-redis/v9"
)

type Repository[T any] struct {
	client *redis.Client
}

func NewCacheRepository[T any](redisClient *redis.Client) *Repository[T] {
	return &Repository[T]{
		client: redisClient,
	}
}

func (repo *Repository[T]) Save(ctx context.Context, key string, value T, ttl time.Duration) error {
	enc := codec.New[T]()
	val, err := enc.Encode(value)
	if err != nil {
		return err
	}

	r := repo.client.Set(ctx, key, val, ttl)
	_, err = r.Result()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository[T]) Get(ctx context.Context, key string) (string, error) {
	return repo.client.Get(ctx, key).Result()
}
