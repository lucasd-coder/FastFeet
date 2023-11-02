package auth

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) CreateUser(ctx context.Context, pld *Register) (*pb.RegisterResponse, error) {
	if err := pld.Validate(s.validate); err != nil {
		return nil, shared.ValidationErrors(err)
	}

	log := logger.FromContext(ctx)

	log.Info("called keycloak")

	id, err := s.repository.Register(ctx, pld)
	if err != nil {
		log.Error(err)
		return nil, shared.CheckError(err)
	}

	return &pb.RegisterResponse{
		Id: id,
	}, nil
}
