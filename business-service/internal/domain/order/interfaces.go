package order

import (
	"context"

	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
)

type (
	Repository interface {
		Save(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error)
		GetAllOrder(ctx context.Context,
			req *pb.GetOrderServiceAllOrderRequest) (*pb.GetAllOrderResponse, error)
	}

	Service interface {
		GetAllOrder(ctx context.Context, pld *GetAllOrderRequest) (*pb.GetAllOrderResponse, error)
		CreateOrder(ctx context.Context, pld Payload) (*pb.OrderResponse, error)
	}
)
