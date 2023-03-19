package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/provider/managerservice"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
)

type UserRepository struct {
	cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{cfg}
}

func (r *UserRepository) Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(r.cfg)
	if err != nil {
		log.Errorf("integration user-manager-service Error: %+v", err)
		return nil, fmt.Errorf("integration user-manager-service Error: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.Save(ctx, req)

	return resp, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(r.cfg)
	if err != nil {
		log.Errorf("integration user-manager-service Error: %+v", err)
		return nil, fmt.Errorf("integration user-manager-service Error: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.FindByEmail(ctx, req)

	return resp, err
}

func (r *UserRepository) FindByCpf(ctx context.Context, req *pb.UserByCpfRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(r.cfg)
	if err != nil {
		log.Errorf("integration user-manager-service Error: %+v", err)
		return nil, fmt.Errorf("integration user-manager-service Error: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.FindByCpf(ctx, req)

	return resp, err
}
