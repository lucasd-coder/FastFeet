package user

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(Service), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	validate   shared.Validator
	repository Repository
}

func NewService(val shared.Validator, repo Repository) *ServiceImpl {
	return &ServiceImpl{
		validate:   val,
		repository: repo,
	}
}
