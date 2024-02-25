package auth

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
)

type (
	Service interface {
		CreateUser(ctx context.Context, pld *Register) (*pb.RegisterResponse, error)
		FindUserByEmail(ctx context.Context, pld *FindUserByEmail) (*pb.GetUserResponse, error)
		GetRoles(ctx context.Context, pld *GetUserID) (*pb.GetRolesResponse, error)
		IsActiveUser(ctx context.Context, pld *GetUserID) (*pb.IsActiveUserResponse, error)
	}

	Repository interface {
		Register(ctx context.Context, pld *Register) (string, error)
		FindUserByEmail(ctx context.Context, pld *FindUserByEmail) (*UserRepresentation, error)
		GetRoles(ctx context.Context, pld *GetUserID) ([]string, error)
		IsActiveUser(ctx context.Context, pld *GetUserID) (bool, error)
	}
)
