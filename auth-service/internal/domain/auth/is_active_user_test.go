package auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IsActiveUserSuite struct {
	suite.Suite
	svc  auth.Service
	repo *mocks.Repository_internal_domain_auth
	ctx  context.Context
}

func (suite *IsActiveUserSuite) SetupTest() {
	val := validator.NewValidation()
	repo := new(mocks.Repository_internal_domain_auth)

	suite.repo = repo
	suite.svc = auth.NewService(val, repo)
	suite.ctx = context.Background()
}

func (suite *IsActiveUserSuite) TestIsActiveUserValidateFailure() {
	pld := &auth.GetUserID{}
	_, err := suite.svc.IsActiveUser(suite.ctx, pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.InvalidArgument)
}

func (suite *IsActiveUserSuite) TestIsActiveUserFailure() {
	pld := &auth.GetUserID{
		ID: "433c311b-93a5-45c3-99c9-b52f3c4aef4f",
	}
	errUserNotFound := fmt.Errorf("fail called FindUserByEmail %w", shared.ErrUserNotFound)

	suite.repo.On("IsActiveUser", suite.ctx, pld).Return(false, errUserNotFound)

	_, err := suite.svc.IsActiveUser(suite.ctx, pld)
	suite.NotNil(err)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *IsActiveUserSuite) TestIsActiveUserSuccess() {
	pld := &auth.GetUserID{
		ID: "433c311b-93a5-45c3-99c9-b52f3c4aef4f",
	}

	suite.repo.On("IsActiveUser", suite.ctx, pld).Return(true, nil)
	resp, err := suite.svc.IsActiveUser(suite.ctx, pld)
	suite.Nil(err)
	suite.Equal(resp.GetActive(), true)
}

func TestIsActiveUserSuite(t *testing.T) {
	suite.Run(t, new(IsActiveUserSuite))
}
