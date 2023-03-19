package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userRepository model.UserRepository
	cfg            *config.Config
}

func NewUserHandler(repo model.UserRepository, cfg *config.Config) *UserHandler {
	return &UserHandler{userRepository: repo, cfg: cfg}
}

func (h *UserHandler) Handler(ctx context.Context, m []byte) error {
	var pld model.User
	if err := json.Unmarshal(m, &pld); err != nil {
		return fmt.Errorf("err Unmarshal: %w", err)
	}
	return h.create(ctx, &pld)
}

func (h *UserHandler) create(ctx context.Context, pld *model.User) error {
	log := logger.FromContext(ctx)

	log.WithFields(logrus.Fields{
		"payload": pld,
	}).Info("received payload")

	if err := pld.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := h.validadeUserWithEmail(ctx, pld.Email); err != nil {
		log.Error("error when validating the email")
		return err
	}

	if err := h.validadeUserWithCpf(ctx, pld.CPF); err != nil {
		log.Error("error when validating the cpf")
		return err
	}

	req := &pb.UserRequest{
		Id:         pld.ID,
		Name:       pld.Name,
		Email:      pld.Email,
		Cpf:        pld.CPF,
		Attributes: pld.Attributes,
	}

	user, err := h.userRepository.Save(ctx, req)
	if err != nil {
		return fmt.Errorf("error when calling save: %w", err)
	}

	log.Infof("payload successfully processed for id: %s", user.Id)

	return nil
}

func (h *UserHandler) validadeUserWithEmail(ctx context.Context, email string) error {
	log := logger.FromContext(ctx)
	userByEmailRequest := &pb.UserByEmailRequest{
		Email: email,
	}

	user, err := h.userRepository.FindByEmail(ctx, userByEmailRequest)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			log.Errorf("fail finByEmail: %v", err)
			return err
		}
	}

	if user != nil {
		log.Errorf("already exist user with id: %s", user.Id)
		return nil
	}

	return nil
}

func (h *UserHandler) validadeUserWithCpf(ctx context.Context, cpf string) error {
	log := logger.FromContext(ctx)

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: cpf,
	}

	user, err := h.userRepository.FindByCpf(ctx, userByCpfRequest)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			log.Errorf("fail findByCpf: %v", err)
			return err
		}
	}

	if user != nil {
		log.Errorf("already exist user with id: %s", user.Id)
		return nil
	}
	return nil
}
