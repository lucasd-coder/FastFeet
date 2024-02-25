//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/pkg/mongodb"
	"github.com/lucasd-coder/fast-feet/user-manger-service/config"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/domain/user/repository"
	val "github.com/lucasd-coder/fast-feet/user-manger-service/internal/provider/validator"
)

func InitializeUserRepository() *repository.UserRepository {
	wire.Build(config.GetConfig, mongodb.GetClientMongoDB, repository.NewUserRepository)
	return &repository.UserRepository{}
}

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return nil
}
