package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/cache"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	ctx         context.Context
	redisServer *miniredis.Miniredis
	repository  *cache.Repository[string]
}

func (suite *RepositorySuite) SetupSuite() {
	suite.ctx = context.Background()
	var err error
	suite.redisServer, err = miniredis.Run()
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *RepositorySuite) SetupTest() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: suite.redisServer.Addr(),
	})
	repo := cache.NewCacheRepository[string](redisClient)
	suite.repository = repo
}

func (suite *RepositorySuite) TearDownSuite() {
	suite.T().Cleanup(func() {
		suite.redisServer.Close()
	})
}

func (suite *RepositorySuite) TestSave() {
	ctx := suite.ctx

	err := suite.repository.Save(ctx, "data", "something here", time.Minute)
	suite.NoError(err)
}

func (suite *RepositorySuite) TestGet() {
	ctx := suite.ctx

	err := suite.repository.Save(ctx, "data", "something here", time.Minute)
	suite.NoError(err)

	getResult, err := suite.repository.Get(ctx, "data")

	var cached string

	enc := codec.New[string]()
	if err := enc.Decode([]byte(getResult), &cached); err != nil {
		suite.Error(err)
	}

	suite.Nil(err)
	suite.Equal(cached, "something here")
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
