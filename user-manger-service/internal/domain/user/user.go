package model

import (
	"time"

	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (user *User) Validate(val shared.Validator) error {
	return val.ValidateStruct(user)
}

func (user *User) GetCreatedAt() string {
	return user.CreatedAt.Format(time.RFC3339)
}
