package order

import (
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

type Payload struct {
	Data      Data   `json:"data,omitempty" validate:"required"`
	EventDate string `json:"eventDate,omitempty" validate:"required,rfc3339"`
}

type Data struct {
	UserID        string  `json:"userId,omitempty" validate:"required,uuid4"`
	DeliverymanID string  `json:"deliverymanId,omitempty" validate:"required,uuid4" `
	Product       Product `json:"product,omitempty" validate:"required"`
	Address       Address `json:"address,omitempty" validate:"required"`
}

type Product struct {
	Name string `json:"name,omitempty" validate:"required,pattern"`
}

type Address struct {
	PostalCode string `json:"postalCode,omitempty" validate:"min=8,max=8,pattern"`
	Number     int32  `json:"number,omitempty" validate:"numeric=integer"`
}

func (payload *Payload) Validate(val shared.Validator) error {
	return val.ValidateStruct(payload)
}

type GetAllOrderRequest struct {
	ID            string     `json:"id,omitempty" validate:"objectID"`
	UserID        string     `json:"userId,omitempty" validate:"required,uuid4"`
	DeliverymanID string     `json:"deliverymanId,omitempty" validate:"required,uuid4"`
	StartDate     string     `json:"startDate,omitempty" validate:"rfc3339"`
	EndDate       string     `json:"endDate,omitempty" validate:"rfc3339"`
	CreatedAt     string     `json:"createdAt,omitempty" validate:"rfc3339"`
	UpdatedAt     string     `json:"updatedAt,omitempty" validate:"rfc3339"`
	CanceledAt    string     `json:"canceledAt,omitempty" validate:"rfc3339"`
	Limit         int64      `json:"limit,omitempty" validate:"numeric=integer"`
	Offset        int64      `json:"offset,omitempty" validate:"numeric=integer"`
	Product       GetProduct `json:"product,omitempty" validate:"required"`
	Address       GetAddress `json:"addresses,omitempty" validate:"required"`
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

func (g *GetAllOrderRequest) Validate(val shared.Validator) error {
	return val.ValidateStruct(g)
}
