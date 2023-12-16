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

type GetRolesSuite struct {
	suite.Suite
	svc  auth.Service
	repo *mocks.Repository_internal_domain_auth
	ctx  context.Context
}

func (suite *GetRolesSuite) SetupTest() {
	val := validator.NewValidation()
	repo := new(mocks.Repository_internal_domain_auth)

	suite.repo = repo
	suite.svc = auth.NewService(val, repo)
	suite.ctx = context.Background()
}

func (suite *GetRolesSuite) TestGetRolesValidateFailure() {
	pld := auth.GetUserID{}
	_, err := suite.svc.GetRoles(suite.ctx, &pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.InvalidArgument)
}

func (suite *GetRolesSuite) TestGetRoleFailure() {
	pld := &auth.GetUserID{
		ID: "433c311b-93a5-45c3-99c9-b52f3c4aef4f",
	}
	errUserNotFound := fmt.Errorf("fail called FindUserByEmail %w", shared.ErrUserNotFound)

	suite.repo.On("GetRoles", suite.ctx, pld).Return(nil, errUserNotFound)

	_, err := suite.svc.GetRoles(suite.ctx, pld)
	suite.NotNil(err)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *GetRolesSuite) TestGetRolesSuccess() {
	pld := &auth.GetUserID{
		ID: "433c311b-93a5-45c3-99c9-b52f3c4aef4f",
	}

	roles := []string{"user"}

	suite.repo.On("GetRoles", suite.ctx, pld).Return(roles, nil)
	resp, err := suite.svc.GetRoles(suite.ctx, pld)
	suite.Nil(err)
	suite.Equal(resp.GetRoles(), roles)
}

func TestGetRolesSuite(t *testing.T) {
	suite.Run(t, new(GetRolesSuite))
}
