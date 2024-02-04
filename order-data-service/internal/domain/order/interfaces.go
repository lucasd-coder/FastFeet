package order

import "context"

type (
	OrderRepository interface {
		Save(ctx context.Context, order *Order) (*Order, error)
		FindByID(ctx context.Context, id string) (*Order, error)
		FindAll(ctx context.Context, pld *GetAllOrderRequest) ([]Order, error)
	}
)
