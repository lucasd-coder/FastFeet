package handler_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/ciphers"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateUserHandlerSuite struct {
	suite.Suite
	cfg      config.Config
	ctx      context.Context
	handler  *handler.Handler
	repoUser *mocks.UserRepository_internal_domain_user
	repoAuth *mocks.AuthRepository_internal_shared
	valErrs  noProviderVal.ValidationErrors
}

func (suite *CreateUserHandlerSuite) SetupSuite() {
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

func (suite *CreateUserHandlerSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoUser := new(mocks.UserRepository_internal_domain_user)

	suite.repoAuth = repoAuth
	suite.repoUser = repoUser
	svc := user.NewService(repoUser, repoAuth, val)
	suite.handler = handler.NewHandler(svc, &suite.cfg)
}

func (suite *CreateUserHandlerSuite) TestCreateUser_UnmarshalFailure() {
	body := []byte(`{""}`)

	err := suite.handler.CreateUser(suite.ctx, body)
	suite.Error(err)
}

func (suite *CreateUserHandlerSuite) TestCreateUser_ValidateFailure() {
	pld := &user.Payload{
		Data: user.Data{
			Email: "test validate email",
		},
	}
	enc, err := suite.encode(pld)
	suite.NoError(err)

	err = suite.handler.CreateUser(suite.ctx, enc)
	suite.Error(err)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *CreateUserHandlerSuite) TestCreateUser_AuthAlreadyExist() {
	payload := &user.Payload{
		Data: user.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       payload.Data.Name,
		Email:      payload.Data.Email,
		Attributes: payload.Data.Attributes,
	}
	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: payload.Data.CPF,
	}
	userByEmailRequest := &pb.UserByEmailRequest{
		Email: payload.Data.Email,
	}

	getUserResp := &shared.GetUserResponse{
		ID:       "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Email:    payload.Data.Email,
		Username: payload.Data.Email,
		Enabled:  true,
	}

	pld, err := suite.encode(payload)
	suite.NoError(err)

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	suite.repoAuth.On("FindByEmail", suite.ctx, payload.Data.Email).Return(getUserResp, nil)

	suite.repoUser.On("Save", suite.ctx, mock.Anything).Return(userResp, nil)

	err = suite.handler.CreateUser(suite.ctx, pld)
	suite.Nil(err)
}

func (suite *CreateUserHandlerSuite) TestCreateUser() {
	payload := &user.Payload{
		Data: user.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       payload.Data.Name,
		Email:      payload.Data.Email,
		Attributes: payload.Data.Attributes,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: payload.Data.CPF,
	}
	userByEmailRequest := &pb.UserByEmailRequest{
		Email: payload.Data.Email,
	}
	register := &shared.RegisterUserResponse{
		ID: userResp.Id,
	}

	pld, err := suite.encode(payload)
	suite.NoError(err)

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	suite.repoAuth.On("FindByEmail", suite.ctx, payload.Data.Email).Return(&shared.GetUserResponse{}, nil)

	suite.repoAuth.On("Register", suite.ctx, payload.ToRegister()).Return(register, nil)

	suite.repoUser.On("Save", suite.ctx, mock.Anything).Return(userResp, nil)

	err = suite.handler.CreateUser(suite.ctx, pld)
	suite.NoError(err)
}

func (suite *CreateUserHandlerSuite) encode(pld *user.Payload) ([]byte, error) {
	codec := codec.New[user.Payload]()

	enc, err := codec.Encode(*pld)
	if err != nil {
		return nil, err
	}

	result, err := ciphers.Encrypt(ciphers.ExtractKey([]byte(suite.cfg.AesKey)), enc)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func TestCreateUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(CreateUserHandlerSuite))
}
