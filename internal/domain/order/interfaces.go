package model

import (
	"context"

	"github.com/lucasd-coder/router-service/pkg/pb"
)

type (
	OrderService interface {
		Save(ctx context.Context, order *Order) error
		GetAllOrders(ctx context.Context, pld *GetAllOrderPayload) (*pb.GetAllOrderResponse, error)
	}

	BusinessRepository interface {
		GetAllOrders(ctx context.Context, req *pb.GetAllOrderRequest) (*pb.GetAllOrderResponse, error)
	}
)
