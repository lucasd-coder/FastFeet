package model

import (
	"context"

	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/pb"
)

type (
	ViaCepRepository interface {
		GetAddress(ctx context.Context, cep string) (*shared.ViaCepAddressResponse, error)
	}

	OrderDataRepository interface {
		Save(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error)
		GetAllOrders(ctx context.Context,
			req *pb.GetOrderServiceAllOrderRequest) (*pb.GetAllOrderResponse, error)
	}
)
