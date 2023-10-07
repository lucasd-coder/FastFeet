package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/pkg/val"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate *validator.Validate

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"userId,omitempty" validate:"required,pattern,uuid4"`
	Name       string             `bson:"name,omitempty" validate:"required,pattern"`
	Email      string             `bson:"email,omitempty" validate:"required,email,pattern"`
	CPF        string             `bson:"cpf,omitempty" validate:"required,isCPF"`
	Attributes map[string]string  `bson:"attributes,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty"`
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

func (user *User) GetCreatedAt() string {
	return user.CreatedAt.Format(time.RFC3339)
}
