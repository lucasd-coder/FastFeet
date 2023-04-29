package model

import (
	"context"

	"github.com/lucasd-coder/router-service/internal/shared"
)

type (
	UserService interface {
		Save(ctx context.Context, user *User) error
	}

	Publish interface {
		Send(ctx context.Context, msg *shared.Message) error
	}
)
