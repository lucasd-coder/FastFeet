package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/business-service/pkg/val"
)

var validate *validator.Validate

type Payload struct {
	Data      Data      `json:"data,omitempty" validate:"required,dive"`
	EventDate time.Time `json:"eventDate,omitempty" validate:"required"`
}

type Data struct {
	Name       string            `json:"name,omitempty" validate:"required,pattern"`
	Email      string            `json:"email,omitempty" validate:"required,email,pattern"`
	CPF        string            `json:"cpf,omitempty" validate:"required,isCPF"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Password   string            `json:"password,omitempty" validate:"required,pattern"`
	Authority  string            `json:"authority,omitempty" validate:"required,oneof=ADMIN USER"`
}

type Register struct {
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Authority string `json:"authority,omitempty"`
}

type RegisterUserResponse struct {
	ID string `json:"id,omitempty"`
}

type GetUserResponse struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

type GetToken struct {
	AccessToken      string `json:"access_token,omitempty"`
	ExpiresIn        string `json:"expires_in,omitempty"`
	RefreshExpiresIn string `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not_before_policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type User struct {
	UserID     string            `json:"userId,omitempty"`
	Name       string            `json:"name,omitempty"`
	Email      string            `json:"email,omitempty"`
	CPF        string            `json:"cpf,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (payload *Payload) Validate() error {
	validate = validator.New()

	if err := validate.RegisterValidation("pattern", val.Pattern); err != nil {
		return err
	}

	if err := validate.RegisterValidation("isCPF", val.TagIsCPF); err != nil {
		return err
	}

	return validate.Struct(payload)
}

func (payload *Payload) ToRegister() *Register {
	return &Register{
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
