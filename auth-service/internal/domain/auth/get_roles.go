package auth

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) GetRoles(ctx context.Context, pld *GetUserID) (*pb.GetRolesResponse, error) {
	log := logger.FromContext(ctx)

	log.Infof("received request GetRoles on with id: %s", pld.ID)

	if err := pld.Validate(s.validate); err != nil {
		return nil, shared.ValidationErrors(err)
	}

	result, err := s.repository.GetRoles(ctx, pld)
	if err != nil {
		return nil, shared.CheckError(err)
	}

	return &pb.GetRolesResponse{
		Roles: result,
	}, nil
}
