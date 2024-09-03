package validator_test

import (
	"testing"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/stretchr/testify/suite"
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
		FirstName string
		LastName  string
		Username  string
		Password  string
		Roles     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "should validate model",
			fields:  fields{"", "test validate email", "", "", ""},
			wantErr: true,
		},
		{
			name:    "should validate field username",
			fields:  fields{"maria", "silva", "test validate email", "123345@#", "USER"},
			wantErr: true,
		},
		{
			name:    "should validate field password",
			fields:  fields{"maria", "maria@gmail.com", "test validate cpf", "12345678", "USER"},
			wantErr: true,
		},
		{
			name:    "should validate with success",
			fields:  fields{"maria", "silva", "maria4@gmail.com", "1234567@#", "USER"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payload := &auth.Register{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				Roles:     tt.fields.Roles,
			}

			if err := suite.val.ValidateStruct(payload); (err != nil) != tt.wantErr {
				suite.T().Errorf("Register.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSuite))
}
