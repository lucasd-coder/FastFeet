package user

import (
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

type Payload struct {
	Data      Data   `json:"data,omitempty" validate:"required"`
	EventDate string `json:"eventDate,omitempty" validate:"required"`
}

type Data struct {
	Name       string            `json:"name,omitempty" validate:"required,pattern"`
	Email      string            `json:"email,omitempty" validate:"required,email,pattern"`
	CPF        string            `json:"cpf,omitempty" validate:"required,isCPF"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Password   string            `json:"password,omitempty" validate:"required,pattern"`
	Authority  string            `json:"authority,omitempty" validate:"required,oneof=ADMIN USER"`
}

type User struct {
	UserID     string            `json:"userId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Email      string            `json:"email,omitempty"`
	CPF        string            `json:"cpf,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type FindByEmailRequest struct {
	Email string `json:"email,omitempty" validate:"required,email,pattern"`
}

func (payload *Payload) Validate(val shared.Validator) error {
	return val.ValidateStruct(payload)
}

func (f *FindByEmailRequest) Validate(val shared.Validator) error {
	return val.ValidateStruct(f)
}

func (payload *Payload) ToRegister() *shared.Register {
	return &shared.Register{
		Name:      payload.Data.Name,
		Username:  payload.Data.Email,
		Password:  payload.Data.Password,
		Authority: payload.Data.Authority,
	}
}

func (payload *Payload) ToUser(userID string) *User {
	return &User{
		UserID:     userID,
		Name:       payload.Data.Name,
		Email:      payload.Data.Email,
		CPF:        payload.Data.CPF,
		Attributes: payload.Data.Attributes,
	}
}
