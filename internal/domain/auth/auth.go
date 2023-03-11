package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Auth struct {
	ID          string    `json:"id,omitempty"`
	Username    string    `json:"username,omitempty" validate:"required,email"`
	Password    string    `json:"password,omitempty" validate:"min:8,containsany=!@#?*"`
	DeliveryMan bool      `json:"delivery_man,omitempty" validate:"default=false,boolean"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (auth *Auth) Validate() error {
	validate = validator.New()
	return validate.Struct(auth)
}
