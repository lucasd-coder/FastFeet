//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/business-service/config"
	order "github.com/lucasd-coder/business-service/internal/domain/order"
	orderHandler "github.com/lucasd-coder/business-service/internal/domain/order/handler"
	orderDataService "github.com/lucasd-coder/business-service/internal/domain/order/service"
	user "github.com/lucasd-coder/business-service/internal/domain/user"
	userHandler "github.com/lucasd-coder/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/business-service/pkg/cache"

	authservice "github.com/lucasd-coder/business-service/internal/provider/authservice/repository"
	managerservice "github.com/lucasd-coder/business-service/internal/provider/managerservice/repository"
	orderdataservice "github.com/lucasd-coder/business-service/internal/provider/orderdataservice/repository"
	val "github.com/lucasd-coder/business-service/internal/provider/validator"
	viacepservice "github.com/lucasd-coder/business-service/internal/provider/viacepservice/repository"
	"github.com/lucasd-coder/business-service/internal/shared"
)

var initializeValidator = wire.NewSet(
	wire.Struct(new(val.Validation)),
	wire.Bind(new(shared.Validator), new(*val.Validation)),
)

var initializeUserRepository = wire.NewSet(
	wire.Bind(new(user.UserRepository), new(*managerservice.UserRepository)),
	managerservice.NewUserRepository,
)

var initializeAuthRepository = wire.NewSet(
	wire.Bind(new(shared.AuthRepository), new(*authservice.AuthRepository)),
	authservice.NewAuthRepository,
)

var initializeViaCepRepository = wire.NewSet(
	wire.Bind(new(order.ViaCepRepository), new(*viacepservice.ViaCepRepository)),
	cache.GetClient,
	viacepservice.NewViaCepRepository,
)

var initializeOrderDataRepository = wire.NewSet(
	wire.Bind(new(order.OrderDataRepository), new(*orderdataservice.OrderDataRepository)),
	orderdataservice.NewOrderDataRepository,
)

func InitializeUserHandler() *userHandler.UserHandler {
	wire.Build(initializeUserRepository,
		initializeAuthRepository, config.GetConfig, initializeValidator, userHandler.NewUserHandler)
	return nil
}

func InitializeOrderHandler() *orderHandler.OrderHandler {
	wire.Build(initializeAuthRepository, initializeViaCepRepository, initializeOrderDataRepository,
		config.GetConfig, initializeValidator, orderHandler.NewOrderHandler)
	return nil
}

func InitializeOrderDataService() *orderDataService.OrderDataService {
	wire.Build(initializeValidator, config.GetConfig, initializeOrderDataRepository, initializeAuthRepository, orderDataService.NewOrderDataService)
	return nil
}
