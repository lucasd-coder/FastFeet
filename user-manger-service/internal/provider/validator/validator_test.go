package validator_test

import (
	"testing"

	model "github.com/lucasd-coder/fast-feet/user-manger-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/shared"
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
		ID         primitive.ObjectID
		Name       string
		Email      string
		CPF        string
		Attributes map[string]string
		Password   string
		Authority  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "should validate model",
			fields:  fields{primitive.ObjectID{}, "", "test validate email", "", map[string]string{}, "", ""},
			wantErr: true,
		},
		{
			name: "should validate field email",
			fields: fields{primitive.ObjectID{}, "maria", "test validate email",
				"901.940.000-28", map[string]string{}, "USER", "12345678"},
			wantErr: true,
		},
		{
			name: "should validate field cpf",
			fields: fields{primitive.ObjectID{}, "maria",
				"maria@gmail.com", "test validate cpf", map[string]string{}, "USER", "12345678"},
			wantErr: true,
		},
		{
			name:    "should validate field password",
			fields:  fields{primitive.ObjectID{}, "maria", "maria2@gmail.com", "995.563.460-07", map[string]string{}, "USER", ""},
			wantErr: true,
		},
		{
			name:    "should validate with success",
			fields:  fields{primitive.ObjectID{}, "maria", "maria4@gmail.com", "999.388.560-63", map[string]string{}, "USER", ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payload := &model.User{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				CPF:        tt.fields.CPF,
				Attributes: tt.fields.Attributes,
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
