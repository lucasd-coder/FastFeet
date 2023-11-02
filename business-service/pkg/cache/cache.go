package cache

import (
	"context"

	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
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

	client = redisClient

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		log.Errorf("Error Redis InstrumentTracing: %v", err)
		return
	}

	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		log.Errorf("Error Redis InstrumentMetrics: %v", err)
		return
	}

	log.Info("Redis Connected")
}

func GetClient() *redis.Client {
	return client
}
