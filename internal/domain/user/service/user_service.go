package service

import (
	"context"
	"fmt"

	model "github.com/lucasd-coder/router-service/internal/domain/user"
	"github.com/lucasd-coder/router-service/internal/provider/validator"
	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type UserService struct {
	validate shared.Validator
}

func NewUserService(validate *validator.Validation) *UserService {
	return &UserService{validate: validate}
}

func (s *UserService) Save(ctx context.Context, user *model.User) error {
	log := logger.FromContext(ctx)

	if err := user.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return msg
	}

	return nil
}
