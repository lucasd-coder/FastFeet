package user

import (
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
)

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

type FindByEmailRequest struct {
	Email string `json:"email,omitempty" validate:"required,email,pattern"`
}

func (user *User) Validate(val shared.Validator) error {
	return val.ValidateStruct(user)
}

func (f *FindByEmailRequest) Validate(val shared.Validator) error {
	return val.ValidateStruct(f)
}
