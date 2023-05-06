//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/order-data-service/internal/domain/delivery/service"
	val "github.com/lucasd-coder/order-data-service/internal/provider/validator"
)

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return &val.Validation{}
}

func InitializeDeliveryService() *service.DeliveryService {
	wire.Build(InitializeValidator, service.NewDeliveryService)
	return &service.DeliveryService{}
}
