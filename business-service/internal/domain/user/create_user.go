package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceImpl) Save(ctx context.Context, pld *Payload) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	if err := pld.Validate(s.validate); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.validadeUserWithEmail(ctx, pld.Data.Email); err != nil {
		log.Errorf("error when validating the email: %v", err)
		return nil, err
	}

	if err := s.validadeUserWithCpf(ctx, pld.Data.CPF); err != nil {
		log.Errorf("error when validating the cpf: %v", err)
		return nil, err
	}

	register, err := s.registerAndReturn(ctx, pld)
	if err != nil {
		log.Errorf("error while call auth-service err: %v", err)
		return nil, err
	}

	req := &pb.UserRequest{
		Id:         register.ID,
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Cpf:        pld.Data.CPF,
		Attributes: pld.Data.Attributes,
	}

	user, err := s.userRepository.Save(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error when calling save: %w", err)
	}

	return user, nil
}

func (s *ServiceImpl) validadeUserWithEmail(ctx context.Context, email string) error {
	log := logger.FromContext(ctx)

	req := &pb.UserByEmailRequest{
		Email: email,
	}
	user, err := s.userRepository.FindByEmail(ctx, req)
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

func (s *ServiceImpl) validadeUserWithCpf(ctx context.Context, cpf string) error {
	log := logger.FromContext(ctx)

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: cpf,
	}

	user, err := s.userRepository.FindByCpf(ctx, userByCpfRequest)
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

func (s *ServiceImpl) registerAndReturn(ctx context.Context, pld *Payload) (*shared.RegisterUserResponse, error) {
	log := logger.FromContext(ctx)

	user, err := s.authRepository.FindByEmail(ctx, pld.Data.Email)
	if err != nil {
		if !errors.Is(err, shared.ErrUserNotFound) {
			log.Errorf("err while call auth-service FindByEmail: %v", err)
			return nil, err
		}
	}

	if user == nil || user.ID == "" {
		register, err := s.authRepository.Register(ctx, pld.ToRegister())
		if err != nil {
			log.Errorf("err while call auth-service Register: %v", err)
			return nil, err
		}
		return register, nil
	}

	return &shared.RegisterUserResponse{
		ID: user.ID,
	}, nil
}
