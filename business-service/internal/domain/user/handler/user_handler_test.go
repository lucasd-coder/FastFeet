package handler_test

import (
	"context"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type UserHandlerSuite struct {
	suite.Suite
	cfg         config.Config
	userHandler *handler.UserHandler
	ctx         context.Context
	conn        *grpc.ClientConn
	repoUser    *mocks.UserRepository_internal_domain_user
	repoAuth    *mocks.AuthRepository_internal_shared
}

func (suite *UserHandlerSuite) SetupSuite() {
	suite.ctx = context.Background()
	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}
	os.Setenv("REDIS_HOST_PASSWORD", "123456")
	os.Setenv("RABBIT_SERVER_URL", "amqp://localhost:5672/fastfeet")

	staticDir := filepath.Join(baseDir, "..", "..", "..", "..", "/config/config-test.yml")

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

func (suite *UserHandlerSuite) SetupTest() {
	lis := bufconn.Listen(1024 * 1024)
	suite.T().Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	suite.T().Cleanup(func() {
		srv.Stop()
	})

	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoUser := new(mocks.UserRepository_internal_domain_user)

	suite.repoAuth = repoAuth
	suite.repoUser = repoUser
	svc := user.NewService(repoUser, repoAuth, val)
	hdler := handler.NewHandler(svc, &suite.cfg)
	suite.userHandler = handler.NewUserHandler(*hdler)

	pb.RegisterUserHandlerServer(srv, suite.userHandler)

	go func() {
		if err := srv.Serve(lis); err != nil {
			suite.T().Error(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 3*time.Second)
	suite.T().Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.T().Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		suite.T().Fatalf("grpc.DialContext %v", err)
	}
	suite.conn = conn
}

func (suite *UserHandlerSuite) TestFindByEmail() {
	client := pb.NewUserHandlerClient(suite.conn)

	in := &pb.UserByEmailRequest{
		Email: "maria12@gmail.com",
	}

	user := &pb.UserResponse{
		Id:    "1c42d3bf-6f10-40b6-94d6-e412c6287e5a",
		Name:  "maria",
		Email: "maria12@gmail.com",
	}

	suite.repoUser.On("FindByEmail", mock.Anything, mock.Anything).Return(user, nil)

	resp, err := client.FindByEmail(suite.ctx, in)
	suite.NoError(err)
	suite.Equal(resp.GetEmail(), user.GetEmail())
	suite.Equal(resp.GetId(), user.GetId())
}

func (suite *UserHandlerSuite) TestFindByEmail_ValidateFailure() {
	client := pb.NewUserHandlerClient(suite.conn)

	in := &pb.UserByEmailRequest{
		Email: "email invalid",
	}

	_, err := client.FindByEmail(suite.ctx, in)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UserHandlerSuite) TestFindByEmail_UserNotFound() {
	client := pb.NewUserHandlerClient(suite.conn)

	in := &pb.UserByEmailRequest{
		Email: "maria12@gmail.com",
	}

	suite.repoUser.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, shared.NotFoundError(shared.ErrUserNotFound))

	_, err := client.FindByEmail(suite.ctx, in)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}
