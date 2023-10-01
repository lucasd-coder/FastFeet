package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/business-service/internal/provider/orderdataservice"

	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
)

type OrderDataRepository struct {
	cfg *config.Config
}

func NewOrderDataRepository(cfg *config.Config) *OrderDataRepository {
	return &OrderDataRepository{cfg: cfg}
}

func (r *OrderDataRepository) Save(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	log := logger.FromContext(ctx)

	conn, err := orderdataservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration save: %+v", err)
		return nil, fmt.Errorf("err while integration save: %w", err)
	}

	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	return client.Save(ctx, req)
}

func (r *OrderDataRepository) GetAllOrder(ctx context.Context,
	req *pb.GetOrderServiceAllOrderRequest) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)
	conn, err := orderdataservice.NewClient(ctx, r.cfg)
	if err != nil {
		log.Errorf("err while integration save: %+v", err)
		return nil, fmt.Errorf("err while integration save: %w", err)
	}

	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	return client.GetAllOrder(ctx, req)
}
