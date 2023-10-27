package model

import (
	"time"

	"github.com/lucasd-coder/fast-feet/order-data-service/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	DeliverymanID string             `bson:"deliverymanId,omitempty" validate:"required,uuid4"`
	CanceledAt    time.Time          `bson:"canceledAt,omitempty"`
	Product       Product            `bson:"product,omitempty" validate:"required"`
	Address       Address            `bson:"addresses,omitempty" validate:"required"`
	SignatureID   string             `bson:"signatureId,omitempty"`
	RecipientID   string             `bson:"recipientId,omitempty"`
	StartDate     time.Time          `bson:"startDate,omitempty" validate:"required"`
	EndDate       time.Time          `bson:"endDate,omitempty" validate:"required"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `bson:"updatedAt,omitempty"`
}

type Product struct {
	Name string `json:"name,omitempty" validate:"required,pattern"`
}

type Address struct {
	Address      string `json:"address,omitempty" validate:"required,pattern"`
	Number       int32  `json:"number,omitempty" validate:"required,numeric=integer"`
	PostalCode   string `json:"postalCode,omitempty" validate:"required,pattern"`
	Neighborhood string `json:"neighborhood,omitempty" validate:"required,pattern"`
	City         string `json:"city,omitempty" validate:"required,pattern"`
	State        string `json:"state,omitempty" validate:"required,pattern"`
}

type CreateOrder struct {
	DeliverymanID string  `json:"deliverymanId,omitempty" validate:"required,uuid4"`
	Product       Product `json:"product,omitempty" validate:"required"`
	Address       Address `json:"addresses,omitempty" validate:"required"`
}

type GetAllOrderRequest struct {
	ID            string     `bson:"_id,omitempty" validate:"pattern"`
	DeliverymanID string     `bson:"deliverymanId,omitempty" validate:"required,uuid4"`
	StartDate     string     `bson:"startDate,omitempty" validate:"rfc3339"`
	EndDate       string     `bson:"endDate,omitempty" validate:"rfc3339"`
	CreatedAt     string     `bson:"createdAt,omitempty" validate:"rfc3339"`
	UpdatedAt     string     `bson:"updatedAt,omitempty" validate:"rfc3339"`
	CanceledAt    string     `bson:"canceledAt,omitempty" validate:"rfc3339"`
	Limit         int64      `bson:"limit,omitempty" validate:"numeric=integer"`
	Offset        int64      `bson:"offset,omitempty" validate:"numeric=integer"`
	Product       GetProduct `bson:"product,omitempty" validate:"required"`
	Address       GetAddress `bson:"addresses,omitempty" validate:"required"`
}

type GetProduct struct {
	Name string `json:"name,omitempty" validate:"pattern"`
}

type GetAddress struct {
	Address      string `bson:"address,omitempty" validate:"pattern"`
	Number       int32  `bson:"number,omitempty" validate:"numeric=integer"`
	PostalCode   string `bson:"postalCode,omitempty" validate:"pattern"`
	Neighborhood string `bson:"neighborhood,omitempty" validate:"pattern"`
	City         string `bson:"city,omitempty" validate:"pattern"`
	State        string `bson:"state,omitempty" validate:"pattern"`
}

func (o *Order) Validate(val shared.Validator) error {
	return val.ValidateStruct(o)
}

func (c *CreateOrder) Validate(val shared.Validator) error {
	return val.ValidateStruct(c)
}

func (g *GetAllOrderRequest) Validate(val shared.Validator) error {
	return val.ValidateStruct(g)
}

func (o *Order) GetCanceledAt() string {
	return o.CanceledAt.Format(time.RFC3339)
}

func (o *Order) GetStartDate() string {
	return o.StartDate.Format(time.RFC3339)
}

func (o *Order) GetEndDate() string {
	return o.EndDate.Format(time.RFC3339)
}

func (o *Order) GetCreatedAt() string {
	return o.CreatedAt.Format(time.RFC3339)
}

func (o *Order) GetUpdatedAt() string {
	return o.UpdatedAt.Format(time.RFC3339)
}

func NewOrder(create CreateOrder) *Order {
	return &Order{
		DeliverymanID: create.DeliverymanID,
		Product:       create.Product,
		Address:       create.Address,
		CreatedAt:     time.Now(),
	}
}

func NewProduct(name string) Product {
	return Product{
		Name: name,
	}
}

func (g *GetAllOrderRequest) GetLimit() int64 {
	if g.Limit == 0 {
		g.Limit = 10
	}

	return g.Limit
}
