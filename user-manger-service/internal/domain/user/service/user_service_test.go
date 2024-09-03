package service_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/user-manger-service/config"
	model "github.com/lucasd-coder/fast-feet/user-manger-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/domain/user/service"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/shared"
	pb "github.com/lucasd-coder/fast-feet/user-manger-service/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	email           = "maria2@gmail.com"
	msgArg          = "suite.svc.Save() = %v, wantErr %v"
	msgUserNotFound = "user not found"
)

type UserServiceSuite struct {
	suite.Suite
	cfg  config.Config
	svc  service.UserService
	ctx  context.Context
	repo *mocks.UserRepository_internal_domain_user
}

func (suite *UserServiceSuite) SetupSuite() {
	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}

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

func (suite *UserServiceSuite) SetupTest() {
	val := validator.NewValidation()
	repo := new(mocks.UserRepository_internal_domain_user)

	suite.repo = repo
	suite.svc = *service.NewUserService(repo, val)
	suite.ctx = context.Background()
}

func (suite *UserServiceSuite) TestSaveValidation() {
	tests := []struct {
		name    string
		args    *pb.UserRequest
		wantErr bool
	}{
		{
			name: "test validation invalid field userID",
			args: &pb.UserRequest{
				UserId: "userID invalid",
			},
			wantErr: true,
		},
		{
			name: "test validation invalid special characters",
			args: &pb.UserRequest{
				UserId: "34dd2b26-7692-48dc-b37d-9445941ed016ee22262f-6d5f-4044-a7d9-e44a196b808c",
				Name:   "maria$%%&%%$#@",
				Email:  "maria@##%%%",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			_, err := suite.svc.Save(suite.ctx, tt.args)
			if err != nil {
				st, ok := status.FromError(err)
				suite.True(ok, msgArg, err, tt.wantErr)
				suite.Equal(st.Code(), codes.InvalidArgument, msgArg, err, tt.wantErr)
			}
		})
	}
}

func (suite *UserServiceSuite) TestSaveAlreadyExistUser() {
	req := &pb.UserRequest{
		UserId: "9a263868-61bc-4e57-b9f8-c6a0a15d2154",
		Name:   "maria2",
		Email:  email,
		Cpf:    "880.910.510-93",
	}

	user := &model.User{
		UserID: req.GetUserId(),
		Name:   req.GetName(),
		Email:  req.GetEmail(),
		CPF:    req.GetCpf(),
	}

	suite.repo.On("FindByUserID", suite.ctx, user.UserID).Return(user, nil)

	_, err := suite.svc.Save(suite.ctx, req)

	st, _ := status.FromError(err)
	msg := fmt.Sprintf("already exist user with userID: %s", user.UserID)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), codes.AlreadyExists, st.Code())
	assert.Equal(suite.T(), msg, st.Message())
}

func (suite *UserServiceSuite) TestSaveMongoErrClientDisconnected() {
	req := &pb.UserRequest{
		UserId: "e06a7169-7df4-4ada-aac3-b673c9713e91",
		Name:   "maria",
		Email:  "maria3@gmail.com",
		Cpf:    "79020873008",
	}

	suite.repo.On("FindByUserID", suite.ctx, req.UserId).Return(nil, mongo.ErrClientDisconnected)

	_, err := suite.svc.Save(suite.ctx, req)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err, mongo.ErrClientDisconnected)
}

func (suite *UserServiceSuite) TestSaveCreatedSuccessfully() {
	req := &pb.UserRequest{
		UserId: "07c837a1-9489-49f3-a038-51a9aff29abe",
		Name:   "maria",
		Email:  "maria4@gmail.com",
		Cpf:    "79020873008",
	}

	user := &model.User{
		UserID: req.GetUserId(),
		Name:   req.GetName(),
		Email:  req.GetEmail(),
		CPF:    req.GetCpf(),
	}

	suite.repo.On("FindByUserID", suite.ctx, req.GetUserId()).Return(nil, nil)
	suite.repo.On("Save", suite.ctx, mock.Anything).Return(user, nil)

	resp, err := suite.svc.Save(suite.ctx, req)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.UserID, req.GetUserId())
	assert.Equal(suite.T(), user.Email, resp.GetEmail())
	assert.Equal(suite.T(), user.Attributes, resp.GetAttributes())
}

func (suite *UserServiceSuite) TestFindByCpfValidation() {
	tests := []struct {
		name    string
		args    *pb.UserByCpfRequest
		wantErr bool
	}{
		{
			name: "test validation invalid cpf",
			args: &pb.UserByCpfRequest{
				Cpf: "invalid cpf",
			},
			wantErr: true,
		},
		{
			name: "test validation badly formatted CPF",
			args: &pb.UserByCpfRequest{
				Cpf: "56304325",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			_, err := suite.svc.FindByCpf(suite.ctx, tt.args)
			if err != nil {
				st, ok := status.FromError(err)
				suite.True(ok, msgArg, err, tt.wantErr)
				suite.Equal(st.Code(), codes.InvalidArgument, msgArg, err, tt.wantErr)
			}
		})
	}
}

func (suite *UserServiceSuite) TestFindByCpfUserNotFond() {
	req := &pb.UserByCpfRequest{
		Cpf: "563.043.250-88",
	}

	suite.repo.On("FindByCpf", suite.ctx, req.Cpf).Return(nil, mongo.ErrNoDocuments)

	_, err := suite.svc.FindByCpf(suite.ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), codes.NotFound, st.Code())
	assert.Equal(suite.T(), msgUserNotFound, st.Message())
}

func (suite *UserServiceSuite) TestFindByCpfGetUserSuccessfully() {
	req := &pb.UserByCpfRequest{
		Cpf: "440.072.470-05",
	}

	userID := "ee22262f-6d5f-4044-a7d9-e44a196b808c"

	user := &model.User{
		UserID: userID,
		Name:   "maria",
		Email:  "maria5@gmail.com",
		CPF:    req.GetCpf(),
	}

	suite.repo.On("FindByCpf", suite.ctx, req.Cpf).Return(user, nil)

	resp, err := suite.svc.FindByCpf(suite.ctx, req)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.UserID, resp.GetUserId())
	assert.Equal(suite.T(), user.Email, resp.GetEmail())
	assert.Equal(suite.T(), user.Attributes, resp.GetAttributes())
}

func (suite *UserServiceSuite) TestFindByEmailValidation() {
	tests := []struct {
		name    string
		args    *pb.UserByEmailRequest
		wantErr bool
	}{
		{
			name: "test validation invalid email",
			args: &pb.UserByEmailRequest{
				Email: "invalid email",
			},
			wantErr: true,
		},
		{
			name: "test validation email not blank",
			args: &pb.UserByEmailRequest{
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			_, err := suite.svc.FindByEmail(suite.ctx, tt.args)
			if err != nil {
				st, ok := status.FromError(err)
				suite.True(ok, "suite.svc.FindByEmail() = %v, wantErr %v", err, tt.wantErr)
				suite.Equal(st.Code(), codes.InvalidArgument, msgArg, err, tt.wantErr)
			}
		})
	}
}

func (suite *UserServiceSuite) TestFindByEmailUserNotFond() {
	req := &pb.UserByEmailRequest{
		Email: "maria@gmail.com",
	}

	suite.repo.On("FindByEmail", suite.ctx, req.Email).Return(nil, mongo.ErrNoDocuments)

	_, err := suite.svc.FindByEmail(suite.ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), codes.NotFound, st.Code())
	assert.Equal(suite.T(), msgUserNotFound, st.Message())
}

func (suite *UserServiceSuite) TestFindByEmailGetUserSuccessfully() {
	req := &pb.UserByEmailRequest{
		Email: "maria@gmail.com",
	}

	userID := "ee22262f-6d5f-4044-a7d9-e44a196b808c"

	user := &model.User{
		UserID: userID,
		Name:   "maria",
		Email:  req.GetEmail(),
		CPF:    "440.072.470-05",
	}

	suite.repo.On("FindByEmail", suite.ctx, req.GetEmail()).Return(user, nil)

	resp, err := suite.svc.FindByEmail(suite.ctx, req)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.UserID, resp.GetUserId())
	assert.Equal(suite.T(), user.Email, resp.GetEmail())
	assert.Equal(suite.T(), user.Attributes, resp.GetAttributes())
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
