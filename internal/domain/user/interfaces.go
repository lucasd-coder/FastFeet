package model

import (
	"context"
)

type (
	UserService interface {
		Save(ctx context.Context, user *User) error
	}
)
