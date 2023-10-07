package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/managerservice"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

type UserRepository struct {
	cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{cfg}
}

func (r *UserRepository) Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration save: %+v", err)
		return nil, fmt.Errorf("err while integration save: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	return client.Save(ctx, req)
}

func (r *UserRepository) FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration findByEmail: %+v", err)
		return nil, fmt.Errorf("err while integration findByEmail: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.FindByEmail(ctx, req)

	return resp, err
}

func (r *UserRepository) FindByCpf(ctx context.Context, req *pb.UserByCpfRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := managerservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration findByCpf: %+v", err)
		return nil, fmt.Errorf("err while integration findByCpf: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	resp, err := client.FindByCpf(ctx, req)

	return resp, err
}
