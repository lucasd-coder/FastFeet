package auth

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) IsActiveUser(ctx context.Context, pld *GetUserID) (*pb.IsActiveUserResponse, error) {
	log := logger.FromContext(ctx)

	log.Infof("received request IsActiveUser on with id: %s", pld.ID)

	if err := pld.Validate(s.validate); err != nil {
		return nil, shared.ValidationErrors(err)
	}

	result, err := s.repository.IsActiveUser(ctx, pld)
	if err != nil {
		return nil, shared.CheckError(err)
	}

	return &pb.IsActiveUserResponse{
		Active: result,
	}, nil
}
