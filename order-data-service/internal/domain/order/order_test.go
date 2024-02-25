package order_test

import (
	"testing"
	"time"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderSuite struct {
	suite.Suite
	val     shared.Validator
	valErrs noProviderVal.ValidationErrors
}

func (suite *OrderSuite) SetupSuite() {
	val := validator.NewValidation()
	suite.val = val
}

func (suite *OrderSuite) TestOrderValidate() {
	tests := []struct {
		name    string
		arg     order.Order
		wantErr bool
	}{
		{
			name: "test validate success",
			arg: order.Order{
				ID:            primitive.NewObjectID(),
				DeliverymanID: "075f0eef-0891-45ad-a3de-d6684c7f390d",
				Product: order.Product{
					Name: "bola",
				},
				Address: order.Address{
					Address:      "rua das marias",
					Number:       10,
					PostalCode:   "2423323252",
					Neighborhood: "casa 3",
					City:         "Rio grande sul",
					State:        "Rio grande sul",
				},
				StartDate: time.Now(),
				EndDate:   time.Now(),
			},
			wantErr: false,
		}, {
			name:    "test validate failed",
			arg:     order.Order{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := tt.arg.Validate(suite.val); (err != nil) != tt.wantErr {
				suite.ErrorAsf(err, &suite.valErrs, "order.Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *OrderSuite) TestCreateOrderValidate() {
	tests := []struct {
		name    string
		arg     order.CreateOrder
		wantErr bool
	}{
		{
			name: "test validate success",
			arg: order.CreateOrder{
				DeliverymanID: "075f0eef-0891-45ad-a3de-d6684c7f390d",
				Product: order.Product{
					Name: "bola",
				},
				Address: order.Address{
					Address:      "rua das marias",
					Number:       10,
					PostalCode:   "2423323252",
					Neighborhood: "casa 3",
					City:         "Rio grande sul",
					State:        "Rio grande sul",
				},
			},
		},
		{
			name:    "test validate failed",
			arg:     order.CreateOrder{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := tt.arg.Validate(suite.val); (err != nil) != tt.wantErr {
				suite.ErrorAsf(err, &suite.valErrs, "createOrder.Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *OrderSuite) TestGetAllOrderRequestValidate() {
	tests := []struct {
		name    string
		arg     order.GetAllOrderRequest
		wantErr bool
	}{
		{
			name: "test validate success",
			arg: order.GetAllOrderRequest{
				ID:            "65b6897300c21bc9e528111b",
				DeliverymanID: "bccef7de-7adf-4699-89c5-d694002bd74e",
				Product: order.GetProduct{
					Name: "bola",
				},
				Address: order.GetAddress{
					Address:      "rua das marias",
					Number:       10,
					PostalCode:   "2423323252",
					Neighborhood: "casa 3",
					City:         "Rio grande sul",
					State:        "Rio grande sul",
				},
			},
			wantErr: false,
		},
		{
			name:    "test validate failed",
			arg:     order.GetAllOrderRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := tt.arg.Validate(suite.val); (err != nil) != tt.wantErr {
				suite.ErrorAsf(err, &suite.valErrs, "getAllOrderRequest.Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *OrderSuite) TestConstructor() {
	createOrder := order.CreateOrder{
		DeliverymanID: "075f0eef-0891-45ad-a3de-d6684c7f390d",
		Product: order.Product{
			Name: "bola",
		},
		Address: order.Address{
			Address:      "rua das marias",
			Number:       10,
			PostalCode:   "2423323252",
			Neighborhood: "casa 3",
			City:         "Rio grande sul",
			State:        "Rio grande sul",
		},
	}
	or := order.NewOrder(createOrder)

	assert.Equal(suite.T(), createOrder.Address.Address, or.Address.Address,
		"not match expected order.Address.Address")
	assert.Equal(suite.T(), createOrder.Address.City, or.Address.City, "not match expected order.Address.City")
	assert.Equal(suite.T(), createOrder.Address.PostalCode, or.Address.PostalCode, "not match expected order.Address.PostalCode")
	assert.Equal(suite.T(), createOrder.Address.Neighborhood,
		or.Address.Neighborhood, "not match expected order.Address.Neighborhood")
	assert.Equal(suite.T(), createOrder.Address.Number, or.Address.Number, "not match expected order.Address.Number")
	assert.Equal(suite.T(), createOrder.Address.State, or.Address.State, "not match expected order.Address.State")
	assert.Equal(suite.T(), createOrder.Product.Name, or.Product.Name, "not match expected order.Product.Name")
	assert.Equal(suite.T(), createOrder.DeliverymanID, or.DeliverymanID, "not match expected order.DeliverymanID")

	name := "bola"

	product := order.NewProduct(name)

	assert.Equal(suite.T(), name, product.Name, "not match expected order.Product.Name")
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
