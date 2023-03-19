//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/provider/managerservice/repository"
)

func InitializeUserRepository() *repository.UserRepository {
	wire.Build(config.GetConfig, repository.NewUserRepository)
	return &repository.UserRepository{}
}
