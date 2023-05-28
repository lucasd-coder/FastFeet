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

type GetAllOrderRequest struct {
	ID            string     `json:"id,omitempty" validate:"objectID"`
	DeliverymanID string     `json:"deliverymanId,omitempty" validate:"required,uuid4"`
	StartDate     string     `json:"startDate,omitempty" validate:"rfc3339"`
	EndDate       string     `json:"endDate,omitempty" validate:"rfc3339"`
	CreatedAt     string     `json:"createdAt,omitempty" validate:"rfc3339"`
	UpdatedAt     string     `json:"updatedAt,omitempty" validate:"rfc3339"`
	CanceledAt    string     `json:"canceledAt,omitempty" validate:"rfc3339"`
	Limit         int64      `json:"limit,omitempty" validate:"numeric=integer"`
	Offset        int64      `json:"offset,omitempty" validate:"numeric=integer"`
	Product       GetProduct `json:"product,omitempty" validate:"required,dive"`
	Address       GetAddress `json:"addresses,omitempty" validate:"required,dive"`
}

type GetProduct struct {
	Name string `json:"name,omitempty" validate:"pattern"`
}

type GetAddress struct {
	Address      string `json:"address,omitempty" validate:"pattern"`
	Number       int32  `json:"number,omitempty" validate:"numeric=integer"`
	PostalCode   string `json:"postalCode,omitempty" validate:"pattern"`
	Neighborhood string `json:"neighborhood,omitempty" validate:"pattern"`
	City         string `json:"city,omitempty" validate:"pattern"`
	State        string `json:"state,omitempty" validate:"pattern"`
}

type GetAllOrderPayload struct {
	GetAllOrderRequest
	UserID string `json:"userId,omitempty" validate:"required,uuid4"`
}

func (g *GetAllOrderPayload) Validate(val shared.Validator) error {
	return val.ValidateStruct(g)
}
