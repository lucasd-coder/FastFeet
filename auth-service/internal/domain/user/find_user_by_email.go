package user

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) FindUserByEmail(ctx context.Context, pld *FindUserByEmail) (*pb.GetUserResponse, error) {
	log := logger.FromContext(ctx)

	log.Info("received request service findUserByEmail")

	if err := pld.Validate(s.validate); err != nil {
		return nil, shared.ValidationErrors(err)
	}

	result, err := s.repository.FindUserByEmail(ctx, pld)
	if err != nil {
		return nil, shared.CheckError(err)
	}

	return &pb.GetUserResponse{
		Id:       result.ID,
		Username: result.Username,
		Enabled:  result.Enabled,
		Email:    result.Email,
	}, nil
}
