package user

import (
	"context"

	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) FindByEmail(ctx context.Context, pld *FindByEmailRequest) (*pb.UserResponse, error) {
	if err := pld.Validate(s.validate); err != nil {
		return nil, err
	}
	log := logger.FromContext(ctx)

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Email,
	}

	log.Info("calling userRepository")

	user, err := s.userRepository.FindByEmail(ctx, userByEmailRequest)
	if err != nil {
		return nil, err
	}

	return user, nil
}
