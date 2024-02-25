//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/order-data-service/config"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order/repository"
	val "github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/pkg/mongodb"
)

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return &val.Validation{}
}

func InitializeOrderRepository() *repository.OrderRepository {
	wire.Build(config.GetConfig, mongodb.GetClientMongoDB, repository.NewOrderRepository)
	return &repository.OrderRepository{}
}
