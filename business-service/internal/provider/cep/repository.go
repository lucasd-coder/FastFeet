package cep

import (
	"context"

	"github.com/lucasd-coder/fast-feet/business-service/config"
	cacheProvider "github.com/lucasd-coder/fast-feet/business-service/internal/provider/cache"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/redis/go-redis/v9"
)

const (
	spanErrRequest         = "Request Error"
	spanErrResponseStatus  = "Response Status Error"
	spanErrExtractResponse = "Error Extract Response"
)

type Repository interface {
	GetAddress(ctx context.Context, cep string) (*shared.AddressResponse, error)
}

func NewBrasilAbertoRepository(cfg *config.Config,
	redisClient *redis.Client) *BrasilAbertoRepository {
	cacheRepository := cacheProvider.NewCacheRepository[shared.AddressResponse](redisClient)
	return &BrasilAbertoRepository{cfg, cacheRepository}
}

func NewViaCepRepository(cfg *config.Config,
	redisCient *redis.Client) *ViaCepRepository {
	cacheRepository := cacheProvider.NewCacheRepository[shared.AddressResponse](redisCient)
	return &ViaCepRepository{cfg, cacheRepository}
}
