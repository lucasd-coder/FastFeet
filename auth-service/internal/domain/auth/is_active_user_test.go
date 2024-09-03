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

type IsActiveUserSuite struct {
	suite.Suite
	cfg  config.Config
	svc  auth.Service
	repo *mocks.Repository_internal_domain_auth
	ctx  context.Context
}

func (suite *IsActiveUserSuite) SetupSuite() {
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
