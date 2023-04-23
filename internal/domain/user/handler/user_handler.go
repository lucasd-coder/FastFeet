package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/internal/shared/ciphers"
	"github.com/lucasd-coder/business-service/internal/shared/codec"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userRepository model.UserRepository
	authRepository model.AuthRepository
	cfg            *config.Config
}

func NewUserHandler(userRepo model.UserRepository,
	authRepo model.AuthRepository, cfg *config.Config,
) *UserHandler {
	return &UserHandler{
		userRepository: userRepo,
		authRepository: authRepo,
		cfg:            cfg,
	}
}

func (h *UserHandler) Handler(ctx context.Context, m []byte) error {
	var pld model.Payload

	dec, err := ciphers.Decrypt(ciphers.ExtractKey([]byte(h.cfg.AesKey)), m)
	if err != nil {
		return fmt.Errorf("err Decrypt: %w", err)
	}

	codec := codec.New[model.Payload]()

	if err := codec.Decode(dec, &pld); err != nil {
		return fmt.Errorf("err Decode: %w", err)
	}
	return h.create(ctx, &pld)
}

func (h *UserHandler) create(ctx context.Context, pld *model.Payload) error {
	log := logger.FromContext(ctx)

	fields := map[string]interface{}{
		"payload": map[string]string{
			"name": pld.Data.Name,
		},
	}
	log.WithFields(fields).Info("received payload")

	if err := pld.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := h.validadeUserWithEmail(ctx, pld.Data.Email); err != nil {
		log.Errorf("error when validating the email: %v", err)
		return err
	}

	if err := h.validadeUserWithCpf(ctx, pld.Data.CPF); err != nil {
		log.Errorf("error when validating the cpf: %v", err)
		return err
	}

	register, err := h.registerAndReturn(ctx, pld)
	if err != nil {
		log.Errorf("err while call auth-service: %v", err)
		return err
	}

	req := &pb.UserRequest{
		Id:         register.ID,
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Cpf:        pld.Data.CPF,
		Attributes: pld.Data.Attributes,
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

	if user != nil && user.Id != "" {
		return fmt.Errorf("error validating user with email: %w", shared.ErrUserAlreadyExist)
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

	if user != nil && user.Id != "" {
		return fmt.Errorf("error validating user with cpf: %w", shared.ErrUserAlreadyExist)
	}
	return nil
}

func (h *UserHandler) registerAndReturn(ctx context.Context, pld *model.Payload) (*model.RegisterUserResponse, error) {
	log := logger.FromContext(ctx)

	user, err := h.authRepository.FindByEmail(ctx, pld.Data.Email)
	if err != nil {
		if !errors.Is(err, shared.ErrUserNotFound) {
			log.Errorf("err while call auth-service FindByEmail: %v", err)
			return nil, err
		}
	}

	if user == nil || user.ID == "" {
		register, err := h.authRepository.Register(ctx, pld.ToRegister())
		if err != nil {
			log.Errorf("err while call auth-service Register: %v", err)
			return nil, err
		}
		return register, nil
	}

	return &model.RegisterUserResponse{
		ID: user.ID,
	}, nil
}
