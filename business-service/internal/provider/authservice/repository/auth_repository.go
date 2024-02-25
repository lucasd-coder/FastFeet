package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/authservice"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type AuthRepository struct {
	cfg *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	return &AuthRepository{
		cfg: cfg,
	}
}

func (r *AuthRepository) Register(ctx context.Context,
	pld *shared.Register) (*shared.RegisterUserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := authservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration register: %+v", err)
		return nil, fmt.Errorf("err while integration register: %w", err)
	}
	defer conn.Close()

	client := pb.NewRegisterHandlerClient(conn)
	req := buildRegisterRequest(pld)
	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("err while integration createUser: %w", err)
	}

	return buildRegisterUserResponse(resp), nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*shared.GetUserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := authservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration findByEmail: %+v", err)
		return nil, fmt.Errorf("err while integration findByEmail: %w", err)
	}
	defer conn.Close()

	client := pb.NewAuthHandlerClient(conn)

	in := &pb.EmptyRequest{}
	header := metadata.New(map[string]string{"email": email})
	ctx = metadata.NewOutgoingContext(ctx, header)

	resp, err := client.FindUserByEmail(ctx, in)
	if err != nil {
		return nil, err
	}

	return &shared.GetUserResponse{
		ID:       resp.GetId(),
		Email:    resp.GetEmail(),
		Username: resp.GetEmail(),
		Enabled:  resp.GetEnabled(),
	}, nil
}

func (r *AuthRepository) FindRolesByID(ctx context.Context, id string) (*shared.GetRolesResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := authservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration findRolesByID: %+v", err)
		return nil, fmt.Errorf("err while integration findRolesByID: %w", err)
	}
	defer conn.Close()

	client := pb.NewAuthHandlerClient(conn)

	in := &pb.EmptyRequest{}
	header := metadata.New(map[string]string{"id": id})
	ctx = metadata.NewOutgoingContext(ctx, header)

	resp, err := client.GetRoles(ctx, in)
	if err != nil {
		return nil, err
	}

	return &shared.GetRolesResponse{
		Roles: resp.GetRoles(),
	}, nil
}

func (r *AuthRepository) IsActiveUser(ctx context.Context, id string) (*shared.IsActiveUser, error) {
	log := logger.FromContext(ctx)

	conn, err := authservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration isActiveUser: %+v", err)
		return nil, fmt.Errorf("err while integration isActiveUser: %w", err)
	}
	defer conn.Close()

	client := pb.NewAuthHandlerClient(conn)

	in := &pb.EmptyRequest{}
	header := metadata.New(map[string]string{"id": id})
	ctx = metadata.NewOutgoingContext(ctx, header)

	resp, err := client.IsActiveUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return &shared.IsActiveUser{
		Active: resp.GetActive(),
	}, nil
}

func buildRegisterUserResponse(resp *pb.RegisterResponse) *shared.RegisterUserResponse {
	return &shared.RegisterUserResponse{
		ID: resp.GetId(),
	}
}

func buildRegisterRequest(pld *shared.Register) *pb.RegisterRequest {
	return &pb.RegisterRequest{
		FirstName: pld.Name,
		Username:  pld.Username,
		Password:  pld.Password,
		Authority: pb.Roles(pb.Roles_value[pld.Authority]),
	}
}
