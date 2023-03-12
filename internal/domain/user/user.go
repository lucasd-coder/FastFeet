package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/business-service/pkg/val"
)

var validate *validator.Validate

type User struct {
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty" validate:"required,pattern"`
	Email      string            `json:"email,omitempty" validate:"required,email,pattern"`
	CPF        string            `json:"cpf,omitempty" validate:"required,isCPF"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (user *User) Validate() error {
	validate = validator.New()

	if err := validate.RegisterValidation("pattern", val.Pattern); err != nil {
		return err
	}

	if err := validate.RegisterValidation("isCPF", val.TagIsCPF); err != nil {
		return err
	}

	return validate.Struct(user)
}
