package user

import (
	"context"

	"github.com/lucasd-coder/fast-feet/router-service/pkg/pb"
)

type (
	Service interface {
		Save(ctx context.Context, user *User) error
		FindUserByEmail(ctx context.Context, pld *FindByEmailRequest) (*pb.UserResponse, error)
	}
)
