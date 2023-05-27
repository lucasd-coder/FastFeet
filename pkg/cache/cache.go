package cache

import (
	"context"

	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func SetUpRedis(ctx context.Context, cfg *config.Config) {
	log := logger.FromContext(ctx)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		DB:       cfg.RedisDB,
		Password: cfg.RedisPassword,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Errorf("Error Redis connection: %+v", err.Error())
		return
	}

	log.Info("Redis Connected")

	client = redisClient
}

func GetClient() *redis.Client {
	return client
}
