package order

import (
	"context"

	"github.com/lucasd-coder/fast-feet/router-service/pkg/pb"
)

type (
	Service interface {
		Save(ctx context.Context, order *Order) error
		GetAllOrder(ctx context.Context, pld *GetAllOrderPayload) (*pb.GetAllOrderResponse, error)
	}
)
