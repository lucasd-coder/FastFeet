package validator_test

import (
	"testing"

	model "github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/shared"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidatorSuite struct {
	suite.Suite
	val shared.Validator
}

func (suite *ValidatorSuite) SetupSuite() {
	val := validator.NewValidation()
	suite.val = val
}

func (suite *ValidatorSuite) TestValidator() {
	type fields struct {
		ID            primitive.ObjectID
		DeliverymanID string
		ProductName   string
		Address       struct {
			Address      string
			Number       int32
			PostalCode   string
			Neighborhood string
			City         string
			State        string
		}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should validate model",
			fields: fields{primitive.ObjectID{}, "", "bola", struct {
				Address      string
				Number       int32
				PostalCode   string
				Neighborhood string
				City         string
				State        string
			}{}},
			wantErr: true,
		},
		{
			name: "should validate field deliverymanId",
			fields: fields{primitive.ObjectID{}, "test validate deliverymanId", "bola",
				struct {
					Address      string
					Number       int32
					PostalCode   string
					Neighborhood string
					City         string
					State        string
				}{}},
			wantErr: true,
		},
		{
			name: "should validate field productName",
			fields: fields{primitive.ObjectID{}, "cd2b4d50-963f-4b47-9e9f-ca4de1d004be",
				"", struct {
					Address      string
					Number       int32
					PostalCode   string
					Neighborhood string
					City         string
					State        string
				}{
					Address:      "rua das marias",
					Number:       10,
					PostalCode:   "65910-809",
					Neighborhood: "Jardim Tropical",
					City:         "Imperatriz",
					State:        "MA",
				}},
			wantErr: true,
		},
		{
			name: "should validate field address",
			fields: fields{primitive.ObjectID{}, "cd2b4d50-963f-4b47-9e9f-ca4de1d004be", "bola", struct {
				Address      string
				Number       int32
				PostalCode   string
				Neighborhood string
				City         string
				State        string
			}{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payload := &model.Order{
				ID:            tt.fields.ID,
				DeliverymanID: tt.fields.DeliverymanID,
				Product:       model.NewProduct(tt.fields.ProductName),
				Address: model.Address{
					Address:      tt.fields.Address.Address,
					Number:       tt.fields.Address.Number,
					PostalCode:   tt.fields.Address.PostalCode,
					Neighborhood: tt.fields.Address.Neighborhood,
					City:         tt.fields.Address.City,
					State:        tt.fields.Address.State,
				},
			}

			if err := suite.val.ValidateStruct(payload); (err != nil) != tt.wantErr {
				suite.T().Errorf("Payload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSuite))
}
