package user

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/router-service/internal/shared/errors"
	"github.com/lucasd-coder/router-service/pkg/logger"
	"github.com/lucasd-coder/router-service/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceImpl) FindUserByEmail(ctx context.Context, pld *FindByEmailRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	if err := pld.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return nil, msg
	}

	req := &pb.UserByEmailRequest{
		Email: pld.Email,
	}

	resp, err := s.businessRepo.FindByEmail(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("fail call businessRepository err: %w", err)
	}

	return resp, nil
}
