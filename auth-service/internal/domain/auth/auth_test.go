package auth_test

import (
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/stretchr/testify/suite"
)

type AuthTest struct {
	suite.Suite
	val     shared.Validator
	valErrs noProviderVal.ValidationErrors
}

func (suite *AuthTest) SetupSuite() {
	val := validator.NewValidation()
	suite.val = val
}

func (suite *AuthTest) TestValidate() {
	register := auth.Register{
		FirstName: "FirstName",
		Username:  "Username",
	}

	err := register.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTest))
}
