package user

import (
	"github.com/google/wire"
	"github.com/lucasd-coder/business-service/internal/shared"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(Service), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	userRepository Repository
	authRepository shared.AuthRepository
	validate       shared.Validator
}

func NewService(userRepo Repository,
	authRepo shared.AuthRepository,
	val shared.Validator,
) *ServiceImpl {
	return &ServiceImpl{
		userRepository: userRepo,
		authRepository: authRepo,
		validate:       val,
	}
}
