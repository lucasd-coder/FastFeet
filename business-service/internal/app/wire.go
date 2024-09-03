//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	orderHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
	user "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	userHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/cache"

	authservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/authservice/repository"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/cep"
	managerservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/managerservice/repository"
	orderdataservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/orderdataservice/repository"
	val "github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

var (
	initializeUserRepository = wire.NewSet(
		wire.Bind(new(user.Repository), new(*managerservice.UserRepository)),
		managerservice.NewUserRepository,
	)
)

func InitializeUserHandler() *userHandler.Handler {
	wire.Build(initializeUserRepository,
		InitializeAuthRepository, InitializeValidator, user.InitializeService, config.GetConfig, userHandler.NewHandler)
	return nil
}

func InitializeBrasilAbertoRepository() *cep.BrasilAbertoRepository {
	wire.Build(wire.NewSet(
		wire.Bind(new(cep.Repository), new(*cep.BrasilAbertoRepository)),
		config.GetConfig, cache.GetClient,
		cep.NewBrasilAbertoRepository,
	))
	return nil
}

func InitializeViaCepRepository() *cep.ViaCepRepository {
	wire.Build(wire.NewSet(
		wire.Bind(new(cep.Repository), new(*cep.ViaCepRepository)),
		config.GetConfig, cache.GetClient,
		cep.NewViaCepRepository,
	))
	return nil
}

func InitializeValidator() shared.Validator {
	wire.Build(wire.NewSet(
		wire.Struct(new(val.Validation)),
		wire.Bind(new(shared.Validator), new(*val.Validation)),
	))
	return nil
}

func InitializeAuthRepository() shared.AuthRepository {
	wire.Build(wire.NewSet(
		wire.Bind(new(shared.AuthRepository), new(*authservice.AuthRepository)),
		config.GetConfig,
		authservice.NewAuthRepository,
	))
	return nil
}

func InitializeOrderDataRepository() order.Repository {
	wire.Build(wire.NewSet(
		wire.Bind(new(order.Repository), new(*orderdataservice.OrderDataRepository)),
		config.GetConfig,
		orderdataservice.NewOrderDataRepository,
	))
	return nil
}

func InitializeOrderHandler() *orderHandler.Handler {
	wire.Build(InitializeAuthRepository, config.GetConfig, newCepRepository, InitializeOrderDataRepository, InitializeValidator, order.InitializeService, orderHandler.NewHandler)
	return nil
}

func newCepRepository() cep.Repository {
	cfg := config.GetConfig()
	if *cfg.ViaCepEnabled {
		return InitializeViaCepRepository()
	}
	if *cfg.BrasilAbertoEnabled {
		return InitializeBrasilAbertoRepository()
	}

	return nil
}
