package user_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindByEmailSuite struct {
	cfg config.Config
	suite.Suite
	ctx      context.Context
	repoAuth *mocks.AuthRepository_internal_shared
	repoUser *mocks.Repository_internal_domain_user
	svc      user.Service
}

func (suite *FindByEmailSuite) SetupSuite() {
	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}

	os.Setenv("REDIS_HOST_PASSWORD", "123456")
	os.Setenv("RABBIT_SERVER_URL", "amqp://localhost:5672/fastfeet")

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

func (suite *FindByEmailSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoUser := new(mocks.Repository_internal_domain_user)

	suite.repoAuth = repoAuth
	suite.repoUser = repoUser
	suite.svc = user.NewService(repoUser, repoAuth, val)
	suite.ctx = context.Background()
}

func (suite *FindByEmailSuite) TestFindByEmailWhenUserNotFound() {
	pld := &user.FindByEmailRequest{
		Email: "maria@gmail.com",
	}

	suite.repoUser.On("FindByEmail", suite.ctx, mock.Anything).Return(nil, shared.NotFoundError(shared.ErrUserNotFound))

	_, err := suite.svc.FindByEmail(suite.ctx, pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *FindByEmailSuite) TestFindByEmail() {
	pld := &user.FindByEmailRequest{
		Email: "maria@gmail.com",
	}

	user := &pb.UserResponse{
		Id:        "1c8d463a-8247-4ac5-aef5-012dffd52fc3",
		Name:      "maria",
		Email:     pld.Email,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	suite.repoUser.On("FindByEmail", suite.ctx, mock.Anything).Return(user, nil)

	resp, err := suite.svc.FindByEmail(suite.ctx, pld)
	suite.NoError(err)
	suite.Equal(user, resp)
}

func (suite *FindByEmailSuite) TestFindByEmailValidateFailure() {
	pld := &user.FindByEmailRequest{}

	_, err := suite.svc.FindByEmail(suite.ctx, pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.InvalidArgument)
}

func TestFindByEmailSuite(t *testing.T) {
	suite.Run(t, new(FindByEmailSuite))
}
