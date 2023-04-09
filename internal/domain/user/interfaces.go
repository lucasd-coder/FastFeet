package model

import (
	"context"
)

type (
	UserRepository interface {
		Save(ctx context.Context, user *User) (*User, error)
		FindByEmail(ctx context.Context, email string) (*User, error)
		FindByUserID(ctx context.Context, userID string) (*User, error)
		FindByCpf(ctx context.Context, cpf string) (*User, error)
	}
)
