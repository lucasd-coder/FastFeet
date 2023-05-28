//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/controller"
	order "github.com/lucasd-coder/router-service/internal/domain/order"
	orderServi "github.com/lucasd-coder/router-service/internal/domain/order/service"
	userServi "github.com/lucasd-coder/router-service/internal/domain/user/service"
	businessservice "github.com/lucasd-coder/router-service/internal/provider/businessservice/repository"
	"github.com/lucasd-coder/router-service/internal/provider/publish"
	val "github.com/lucasd-coder/router-service/internal/provider/validator"
	"github.com/lucasd-coder/router-service/internal/shared"
)

func extractOptionOrderEvents() *shared.Options {
	cfg := config.GetConfig()
	return &shared.Options{
		TopicURL:    cfg.TopicOrderEvents.URL,
		MaxRetries:  cfg.TopicOrderEvents.MaxRetries,
		WaitingTime: cfg.TopicOrderEvents.WaitingTime,
	}
}

func extractOptionUserEvents() *shared.Options {
	cfg := config.GetConfig()
	return &shared.Options{
		TopicURL:    cfg.TopicUserEvents.URL,
		MaxRetries:  cfg.TopicUserEvents.MaxRetries,
		WaitingTime: cfg.TopicUserEvents.WaitingTime,
	}
}

func InitializeValidator() *val.Validation {
	wire.Build(val.NewValidation)
	return nil
}

func InitializeOrderEventsPublish() *publish.Published {
	wire.Build(extractOptionOrderEvents, publish.NewPublished)
	return nil
}

func InitializeUserEventsPublish() *publish.Published {
	wire.Build(extractOptionUserEvents, publish.NewPublished)
	return nil
}

var initializeBusinessRepository = wire.NewSet(
	wire.Bind(new(order.BusinessRepository), new(*businessservice.BusinessRepository)),
	businessservice.NewBusinessRepository,
)

func InitializeUserService() *userServi.UserService {
	wire.Build(InitializeValidator, InitializeUserEventsPublish, config.GetConfig, userServi.NewUserService)
	return nil
}

func InitializeUserController() *controller.UserController {
	wire.Build(InitializeUserService, controller.NewUserController)
	return nil
}

func InitializeOrderService() *orderServi.OrderService {
	wire.Build(InitializeValidator, InitializeOrderEventsPublish, config.GetConfig, initializeBusinessRepository, orderServi.NewOrderService)
	return nil
}

func InitializeOrderController() *controller.OrderController {
	wire.Build(InitializeOrderService, controller.NewOrderController)
	return nil
}
