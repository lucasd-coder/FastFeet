// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	handler2 "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"
	repository2 "github.com/lucasd-coder/fast-feet/business-service/internal/provider/authservice/repository"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/cep"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/managerservice/repository"
	repository3 "github.com/lucasd-coder/fast-feet/business-service/internal/provider/orderdataservice/repository"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/cache"
)

import (
	_ "net/http/pprof"
)

// Injectors from wire.go:

func InitializeUserHandler() *handler.Handler {
	configConfig := config.GetConfig()
	userRepository := repository.NewUserRepository(configConfig)
	authRepository := InitializeAuthRepository()
	validator := InitializeValidator()
	serviceImpl := user.NewService(userRepository, authRepository, validator)
	handlerHandler := handler.NewHandler(serviceImpl, configConfig)
	return handlerHandler
}

func InitializeBrasilAbertoRepository() *cep.BrasilAbertoRepository {
	configConfig := config.GetConfig()
	client := cache.GetClient()
	brasilAbertoRepository := cep.NewBrasilAbertoRepository(configConfig, client)
	return brasilAbertoRepository
}

func InitializeViaCepRepository() *cep.ViaCepRepository {
	configConfig := config.GetConfig()
	client := cache.GetClient()
	viaCepRepository := cep.NewViaCepRepository(configConfig, client)
	return viaCepRepository
}

func InitializeValidator() shared.Validator {
	validation := &validator.Validation{}
	return validation
}

func InitializeAuthRepository() shared.AuthRepository {
	configConfig := config.GetConfig()
	authRepository := repository2.NewAuthRepository(configConfig)
	return authRepository
}

func InitializeOrderDataRepository() order.Repository {
	configConfig := config.GetConfig()
	orderDataRepository := repository3.NewOrderDataRepository(configConfig)
	return orderDataRepository
}

func InitializeOrderHandler() *handler2.Handler {
	sharedValidator := InitializeValidator()
	orderRepository := InitializeOrderDataRepository()
	authRepository := InitializeAuthRepository()
	cepRepository := newCepRepository()
	serviceImpl := order.NewService(sharedValidator, orderRepository, authRepository, cepRepository)
	configConfig := config.GetConfig()
	handlerHandler := handler2.NewHandler(serviceImpl, configConfig)
	return handlerHandler
}

// wire.go:

var (
	initializeUserRepository = wire.NewSet(wire.Bind(new(user.Repository), new(*repository.UserRepository)), repository.NewUserRepository)
)

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
