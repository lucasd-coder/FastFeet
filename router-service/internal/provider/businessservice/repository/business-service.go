package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/provider/businessservice"
	"github.com/lucasd-coder/router-service/pkg/logger"
	"github.com/lucasd-coder/router-service/pkg/pb"
)

type BusinessRepository struct {
	cfg *config.Config
}

func NewBusinessRepository(cfg *config.Config) *BusinessRepository {
	return &BusinessRepository{cfg: cfg}
}

func (r *BusinessRepository) GetAllOrder(ctx context.Context, req *pb.GetAllOrderRequest) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := businessservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration save: %+v", err)
		return nil, fmt.Errorf("err while integration save: %w", err)
	}

	defer conn.Close()

	client := pb.NewOrderHandlerClient(conn)

	return client.GetAllOrder(ctx, req)
}

func (r *BusinessRepository) FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := businessservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration save: %+v", err)
		return nil, fmt.Errorf("err while integration save: %w", err)
	}

	defer conn.Close()

	client := pb.NewUserHandlerClient(conn)

	return client.FindByEmail(ctx, req)
}
