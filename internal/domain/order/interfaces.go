package model

import (
	"context"
)

type (
	OrderService interface {
		Save(ctx context.Context, order *Order) error
	}
)
