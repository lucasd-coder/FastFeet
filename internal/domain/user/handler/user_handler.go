package handler

import (
	"context"
	"fmt"

	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userRepository model.UserRepository
}

func NewUserHandler(repo model.UserRepository) *UserHandler {
	return &UserHandler{userRepository: repo}
}

func (h *UserHandler) Create(ctx context.Context, pld *model.User) error {
	log := logger.FromContext(ctx)

	log.WithFields(logrus.Fields{
		"payload": pld,
	}).Info("received payload")

	if err := pld.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Email,
	}

	user, err := h.userRepository.FindByEmail(ctx, userByEmailRequest)
	if err != nil {
		log.Errorf("fail finByEmail: %v", err)
		return err
	}

	log.Info(user)
	return nil
}
