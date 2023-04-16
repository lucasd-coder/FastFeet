//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/business-service/config"
	authservice "github.com/lucasd-coder/business-service/internal/provider/authservice/repository"
	managerservice "github.com/lucasd-coder/business-service/internal/provider/managerservice/repository"
)

func InitializeUserRepository() *managerservice.UserRepository {
	wire.Build(config.GetConfig, managerservice.NewUserRepository)
	return &managerservice.UserRepository{}
}

func InitializeAuthRepository() *authservice.AuthRepository {
	wire.Build(config.GetConfig, authservice.NewAuthRepository)
	return &authservice.AuthRepository{}
}
