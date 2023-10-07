//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/user-manger-service/config"
	"github.com/lucasd-coder/user-manger-service/internal/domain/user/repository"
	"github.com/lucasd-coder/user-manger-service/internal/domain/user/service"
	"github.com/lucasd-coder/user-manger-service/pkg/mongodb"
)

func InitializeUserRepository() *repository.UserRepository {
	wire.Build(config.GetConfig, mongodb.GetClientMongoDB, repository.NewUserRepository)
	return &repository.UserRepository{}
}

func InitializeUserService() *service.UserService {
	wire.Build(InitializeUserRepository, service.NewUserService)
	return &service.UserService{}
}
