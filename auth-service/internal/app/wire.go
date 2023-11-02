//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	authHandler "github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth/handler"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/kecloak"
	val "github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
)

var initializeValidator = wire.NewSet(
	wire.Struct(new(val.Validation)),
	wire.Bind(new(shared.Validator), new(*val.Validation)),
)

var initializeRepository = wire.NewSet(
	wire.Bind(new(auth.Repository), new(*kecloak.Repository)),
	kecloak.NewRepository,
)

func InitializeAuthHandler() *authHandler.Handler {
	wire.Build(initializeValidator, initializeRepository, auth.InitializeService,
		config.GetConfig, authHandler.NewHandler)
	return nil
}
