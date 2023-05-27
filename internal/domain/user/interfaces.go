package model

import (
	"context"

	"github.com/lucasd-coder/business-service/pkg/pb"
)

type (
	UserRepository interface {
		Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error)
		FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error)
		FindByCpf(ctx context.Context, req *pb.UserByCpfRequest) (*pb.UserResponse, error)
	}
)
