package delivery

import (
	"time"

	"github.com/lucasd-coder/order-data-service/internal/shared"
)

type Delivery struct {
	ID            string    `json:"id,omitempty"`
	DeliverymanID string    `json:"deliverymanId,omitempty" validate:"required,pattern"`
	CanceledAt    time.Time `json:"canceledAt,omitempty"`
	Product       Product   `json:"product,omitempty" validate:"required,dive"`
	Address       Address   `json:"addresses,omitempty" validate:"required,dive"`
	SignatureID   string    `json:"signatureId,omitempty"`
	RecipientID   string    `json:"recipientId,omitempty"`
	StartDate     time.Time `json:"startDate,omitempty" validate:"required"`
	EndDate       time.Time `json:"endDate,omitempty" validate:"required"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

type Product struct {
	Name string `json:"name,omitempty" validate:"required,pattern"`
}

type Address struct {
	Address      string `json:"address,omitempty" validate:"required,pattern"`
	Number       int32  `json:"number,omitempty" validate:"required"`
	PostalCode   string `json:"postalCode,omitempty" validate:"required,pattern"`
	Neighborhood string `json:"neighborhood,omitempty" validate:"required,pattern"`
	City         string `json:"city,omitempty" validate:"required,pattern"`
	State        string `json:"state,omitempty" validate:"required,pattern"`
}

func (d *Delivery) Validate(val shared.Validator) error {
	return val.ValidateStruct(d)
}

func (d *Delivery) GetCanceledAt() string {
	return d.CanceledAt.Format(time.RFC3339)
}

func (d *Delivery) GetStartDate() string {
	return d.StartDate.Format(time.RFC3339)
}

func (d *Delivery) GetEndDate() string {
	return d.EndDate.Format(time.RFC3339)
}

func (d *Delivery) GetCreatedAt() string {
	return d.CreatedAt.Format(time.RFC3339)
}

func (d *Delivery) GetUpdatedAt() string {
	return d.UpdatedAt.Format(time.RFC3339)
}

func NewProduct(name string) Product {
	return Product{
		Name: name,
	}
}
