package model

import (
	"context"
)

type (
	UserRepository interface {
		Save(ctx context.Context, user *User) (*User, error)
		FindByEmail(ctx context.Context, email string) (*User, error)
		FindByID(ctx context.Context, id string) (*User, error)
		FindByCpf(ctx context.Context, cpf string) (*User, error)
	}
)
