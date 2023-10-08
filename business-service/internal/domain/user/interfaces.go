package user

import (
	"context"

	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
)

type (
	Repository interface {
		Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error)
		FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error)
		FindByCpf(ctx context.Context, req *pb.UserByCpfRequest) (*pb.UserResponse, error)
	}

	Service interface {
		Save(ctx context.Context, pld *Payload) (*pb.UserResponse, error)
		FindByEmail(ctx context.Context, pld *FindByEmailRequest) (*pb.UserResponse, error)
	}
)
