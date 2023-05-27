package model

import (
	"github.com/lucasd-coder/router-service/internal/shared"
)

type Payload struct {
	Data      Order  `json:"data,omitempty" validate:"required,dive"`
	EventDate string `json:"eventDate,omitempty" validate:"required,rfc3339"`
}

type Order struct {
	UserID        string  `json:"userId,omitempty" validate:"required,uuid4"`
	DeliverymanID string  `json:"deliverymanId,omitempty" validate:"required,uuid4"`
	Product       Product `json:"product,omitempty" validate:"required,dive"`
	Address       Address `json:"address,omitempty" validate:"required,dive"`
}

type Product struct {
	Name string `json:"name,omitempty" validate:"required,pattern"`
}

type Address struct {
	PostalCode string `json:"postalCode,omitempty" validate:"min=8,max=8,pattern"`
	Number     int32  `json:"number,omitempty" validate:"numeric=integer,min=1"`
}

func (order *Order) Validate(val shared.Validator) error {
	return val.ValidateStruct(order)
}

type CreateOrder struct {
	DeliverymanID string  `json:"deliverymanId,omitempty"`
	Product       Product `json:"product,omitempty"`
	Address       Address `json:"address,omitempty"`
}

func (c *CreateOrder) NewOrder(userID string) *Order {
	if c == nil {
		return nil
	}
	return &Order{
		UserID:        userID,
		DeliverymanID: c.DeliverymanID,
		Product:       c.Product,
		Address:       c.Address,
	}
}
