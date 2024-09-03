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
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindUserByEmailSuite struct {
	suite.Suite
	cfg  config.Config
	svc  auth.Service
	repo *mocks.Repository_internal_domain_auth
	ctx  context.Context
}

func (suite *FindUserByEmailSuite) SetupSuite() {
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

func (suite *FindUserByEmailSuite) SetupTest() {
	val := validator.NewValidation()
	repo := new(mocks.Repository_internal_domain_auth)

	suite.repo = repo
	suite.svc = auth.NewService(val, repo)
	suite.ctx = context.Background()
}

func (suite *FindUserByEmailSuite) TestFindUserByEmailValidateFailure() {
	pld := auth.FindUserByEmail{}
	_, err := suite.svc.FindUserByEmail(suite.ctx, &pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.InvalidArgument)
}

func (suite *FindUserByEmailSuite) TestFindUserByEmailFailure() {
	pld := &auth.FindUserByEmail{
		Email: "maria123@gmail.com",
	}

	errUserNotFound := fmt.Errorf("fail called FindUserByEmail %w", shared.ErrUserNotFound)

	suite.repo.On("FindUserByEmail", suite.ctx, pld).Return(nil, errUserNotFound)

	_, err := suite.svc.FindUserByEmail(suite.ctx, pld)
	suite.NotNil(err)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *FindUserByEmailSuite) TestFindUserByEmailSuccess() {
	email := "maria@gmail.com"

	pld := &auth.FindUserByEmail{
		Email: email,
	}

	userRepresentation := &auth.UserRepresentation{
		ID:       "123456",
		Username: "maria",
		Enabled:  true,
		Email:    email,
	}

	rp := &pb.GetUserResponse{
		Id:       "123456",
		Username: "maria",
		Enabled:  true,
		Email:    email,
	}

	suite.repo.On("FindUserByEmail", suite.ctx, pld).Return(userRepresentation, nil)
	resp, err := suite.svc.FindUserByEmail(suite.ctx, pld)
	suite.Nil(err)
	suite.Equal(resp.GetEmail(), rp.GetEmail())
	suite.Equal(resp.GetId(), rp.GetId())
	suite.Equal(resp.GetUsername(), rp.GetUsername())
}

func TestFindUserByEmailSuite(t *testing.T) {
	suite.Run(t, new(FindUserByEmailSuite))
}
