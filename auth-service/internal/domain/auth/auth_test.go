package auth_test

import (
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	val     shared.Validator
	valErrs noProviderVal.ValidationErrors
}

func (suite *AuthSuite) SetupSuite() {
	val := validator.NewValidation()
	suite.val = val
}

func (suite *AuthSuite) TestRegisterValidate() {
	register := auth.Register{
		FirstName: "FirstName",
		Username:  "Username",
	}

	err := register.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *AuthSuite) TestGetRolesEnumTypeAdmin() {
	roles := "ADMIN"

	rolesEnum := auth.GetRolesFromString(roles)
	suite.Equal(rolesEnum.String(), "admin")
}

func (suite *AuthSuite) TestGetRolesEnumTypeUser() {
	roles := "USER"

	rolesEnum := auth.GetRolesFromString(roles)
	suite.Equal(rolesEnum.String(), "user")
}

func (suite *AuthSuite) TestGetRolesEnumTypeUnknown() {
	roles := "not found"

	rolesEnum := auth.GetRolesFromString(roles)
	suite.Equal(rolesEnum.String(), "unknown")
}

func (suite *AuthSuite) TestFindUserByEmailValidate() {
	findUserByEmail := auth.FindUserByEmail{
		Email: "Email",
	}

	err := findUserByEmail.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *AuthSuite) TestGetUserIDValidate() {
	getUserID := auth.GetUserID{
		ID: "12345678",
	}

	err := getUserID.Validate(suite.val)
	suite.ErrorAs(err, &suite.valErrs)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
