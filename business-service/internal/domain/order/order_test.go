package order_test

import (
	"testing"
	"time"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/stretchr/testify/suite"
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

func (suite *OrderSuite) TestPayloadValidate() {
	pld := order.Payload{
		Data:      order.Data{},
		EventDate: "",
	}

	err := pld.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *OrderSuite) TestNewPayload() {
	pld := order.Payload{
		Data: order.Data{
			UserID:        "136784e3-52d1-4a15-b35a-99a817781dc9",
			DeliverymanID: "8ee8b371-7d54-4a90-840d-bfd68f6ad1d2",
			Product: order.Product{
				Name: "batata",
			},
			Address: order.Address{
				PostalCode: "10968676",
				Number:     10,
			},
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	err := pld.Validate(suite.val)
	suite.NoError(err)
}

func (suite *OrderSuite) TestGetAllOrderRequestValidate() {
	pld := order.GetAllOrderRequest{
		UserID:        "",
		DeliverymanID: "",
	}

	err := pld.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *OrderSuite) TestNeWGetAllOrderRequest() {
	pld := order.GetAllOrderRequest{
		UserID:        "4faa122a-4bf6-44c8-ab06-87c511408bf4",
		DeliverymanID: "c77331c3-7c55-443b-b55a-40cbe145a9ab",
	}

	err := pld.Validate(suite.val)
	suite.NoError(err)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
