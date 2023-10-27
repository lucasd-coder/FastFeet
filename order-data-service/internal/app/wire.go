//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/pkg/mongodb"
	"github.com/lucasd-coder/fast-feet/order-data-service/config"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order/repository"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order/service"
	val "github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
)

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return &val.Validation{}
}

func InitializeOrderRepository() *repository.OrderRepository {
	wire.Build(config.GetConfig, mongodb.GetClientMongoDB, repository.NewOrderRepository)
	return &repository.OrderRepository{}
}

func InitializeOrderService() *service.OrderService {
	wire.Build(InitializeValidator, InitializeOrderRepository, service.NewOrderService)
	return &service.OrderService{}
}
