package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/user-manger-service/pkg/val"
)

var validate *validator.Validate

type User struct {
	ID         string
	Name       string `validate:"required,pattern"`
	Email      string `validate:"required,email,pattern"`
	CPF        string `validate:"required,pattern"`
	Attributes map[string]string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (user *User) Validate() error {
	validate = validator.New()

	if err := validate.RegisterValidation("pattern", val.Pattern); err != nil {
		return err
	}

	return validate.Struct(user)
}
