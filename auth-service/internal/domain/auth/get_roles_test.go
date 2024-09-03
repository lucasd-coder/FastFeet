package auth_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRolesSuite struct {
	suite.Suite
	cfg  config.Config
	svc  auth.Service
	repo *mocks.Repository_internal_domain_auth
	ctx  context.Context
}

func (suite *GetRolesSuite) SetupSuite() {
	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}

	staticDir := filepath.Join(baseDir, "..", "..", "..", "/config/config-test.yml")

	slog.Info("config lod", "dir", staticDir)
	err = cleanenv.ReadConfig(staticDir, &suite.cfg)
	if err != nil {
		suite.T().Errorf("cleanenv.ReadConfig() error = %v", err)
		return
	}
	config.ExportConfig(&suite.cfg)

	optlogger := shared.NewOptLogger(&suite.cfg)
	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)
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
