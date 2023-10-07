package shared

import (
	"context"
	"time"
)

type Validator interface {
	ValidateStruct(s interface{}) error
}

type CacheRepository[T any] interface {
	Save(ctx context.Context, key string, value T, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type AuthRepository interface {
	Register(ctx context.Context, pld *Register) (*RegisterUserResponse, error)
	FindByEmail(ctx context.Context, email string) (*GetUserResponse, error)
	FindRolesByID(ctx context.Context, ID string) (*GetRolesResponse, error)
	IsActiveUser(ctx context.Context, ID string) (*IsActiveUser, error)
}
