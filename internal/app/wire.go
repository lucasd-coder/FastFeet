//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/controller"
	"github.com/lucasd-coder/router-service/internal/domain/user/service"
	"github.com/lucasd-coder/router-service/internal/provider/publish"
	val "github.com/lucasd-coder/router-service/internal/provider/validator"
)

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return &val.Validation{}
}

func InitializePublish() *publish.Published {
	wire.Build(config.GetConfig, publish.NewPublished)
	return &publish.Published{}
}

func InitializeUserService() *service.UserService {
	wire.Build(InitializeValidator, InitializePublish, config.GetConfig, service.NewUserService)
	return &service.UserService{}
}

func InitializeUserController() *controller.UserController {
	wire.Build(InitializeUserService, controller.NewUserController)
	return &controller.UserController{}
}
