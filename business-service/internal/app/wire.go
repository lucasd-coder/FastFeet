//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	order "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	orderHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
	user "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	userHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"

	authservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/authservice/repository"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/cep"
	managerservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/managerservice/repository"
	orderdataservice "github.com/lucasd-coder/fast-feet/business-service/internal/provider/orderdataservice/repository"
	val "github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

var initializeValidator = wire.NewSet(
	wire.Struct(new(val.Validation)),
	wire.Bind(new(shared.Validator), new(*val.Validation)),
)

var initializeUserRepository = wire.NewSet(
	wire.Bind(new(user.Repository), new(*managerservice.UserRepository)),
	managerservice.NewUserRepository,
)

var initializeAuthRepository = wire.NewSet(
	wire.Bind(new(shared.AuthRepository), new(*authservice.AuthRepository)),
	authservice.NewAuthRepository,
)

var initializeOrderDataRepository = wire.NewSet(
	wire.Bind(new(order.Repository), new(*orderdataservice.OrderDataRepository)),
	orderdataservice.NewOrderDataRepository,
)

func InitializeUserHandler() *userHandler.Handler {
	wire.Build(initializeUserRepository,
		initializeAuthRepository, initializeValidator, user.InitializeService, config.GetConfig, userHandler.NewHandler)
	return nil
}

func InitializeOrderHandler() *orderHandler.Handler {
	wire.Build(initializeAuthRepository, config.GetConfig, cep.NewRepository, initializeOrderDataRepository, initializeValidator, order.InitializeService, orderHandler.NewHandler)
	return nil
}
