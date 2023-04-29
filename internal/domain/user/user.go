package model

import (
	"github.com/lucasd-coder/router-service/internal/shared"
)

type CreateEvent struct {
	Message string `json:"message,omitempty"`
}

type Payload struct {
	Data      User   `json:"data,omitempty"`
	EventDate string `json:"eventDate,omitempty"`
}
type User struct {
	Name      string `json:"name,omitempty" validate:"required,pattern"`
	Email     string `json:"email,omitempty" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"min=8,containsany=!@#?*"`
	CPF       string `json:"cpf,omitempty" validate:"required,isCPF"`
	Authority string `json:"authority,omitempty" validate:"required,oneof=ADMIN USER"`
}

func (user *User) Validate(val shared.Validator) error {
	return val.ValidateStruct(user)
}
