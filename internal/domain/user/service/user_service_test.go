package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/user-manger-service/config"
	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"github.com/lucasd-coder/user-manger-service/internal/domain/user/service"
	"github.com/lucasd-coder/user-manger-service/internal/mocks"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	pb "github.com/lucasd-coder/user-manger-service/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSave_InvalidObjectID(t *testing.T) {
	req := &pb.UserRequest{
		Id: "id invalid",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.Save(ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSave_InvalidUserRequest(t *testing.T) {
	req := &pb.UserRequest{
		Id:    "6404a984f18d899ec00c2a76",
		Name:  "maria$%%&%%$#@",
		Email: "maria@##%%%",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.Save(ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSave_AlreadyExistUser(t *testing.T) {
	req := &pb.UserRequest{
		Id:    "6404a984f18d899ec00c2a76",
		Name:  "maria",
		Email: "maria@gmail.com",
		Cpf:   "880.910.510-93",
	}

	objectID := objectIDFromHex(req.Id)

	user := &model.User{
		ID:    objectID,
		Name:  req.GetName(),
		Email: req.GetEmail(),
		CPF:   req.GetCpf(),
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	userRepoMock.On("FindByID", ctx, req.Id).Return(user, nil)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.Save(ctx, req)

	st, _ := status.FromError(err)

	msg := fmt.Sprintf("already exist user with id: %s", req.Id)

	assert.NotNil(t, err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Equal(t, msg, st.Message())
}

func TestSave_MongoErrClientDisconnected(t *testing.T) {
	req := &pb.UserRequest{
		Id:    "6404a984f18d899ec00c2a76",
		Name:  "maria",
		Email: "maria@gmail.com",
		Cpf:   "79020873008",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	userRepoMock.On("FindByID", ctx, req.Id).Return(nil, mongo.ErrClientDisconnected)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.Save(ctx, req)

	assert.NotNil(t, err)
	assert.Equal(t, err, mongo.ErrClientDisconnected)
}

func TestSave_CreatedSuccessfully(t *testing.T) {
	req := &pb.UserRequest{
		Id:    "6404a984f18d899ec00c2a76",
		Name:  "maria",
		Email: "maria@gmail.com",
		Cpf:   "79020873008",
	}

	objectID := objectIDFromHex(req.Id)

	user := &model.User{
		ID:    objectID,
		Name:  req.GetName(),
		Email: req.GetEmail(),
		CPF:   req.GetCpf(),
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	userRepoMock.On("FindByID", ctx, req.Id).Return(nil, nil)
	userRepoMock.On("Save", ctx, mock.Anything).Return(user, nil)

	service := service.UserService{UserRepository: userRepoMock}

	resp, err := service.Save(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), resp.GetId())
	assert.Equal(t, user.Email, resp.GetEmail())
	assert.Equal(t, user.Attributes, resp.GetAttributes())
}

func TestFindByCpf_InvalidCpf(t *testing.T) {
	req := &pb.UserByCpfRequest{
		Cpf: "invalid cpf",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.FindByCpf(ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestFindByCpf_UserNotFond(t *testing.T) {
	req := &pb.UserByCpfRequest{
		Cpf: "440.072.470-05",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	userRepoMock.On("FindByCpf", ctx, req.Cpf).Return(nil, mongo.ErrNoDocuments)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.FindByCpf(ctx, req)

	st, _ := status.FromError(err)

	msg := "user not found"

	assert.NotNil(t, err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, msg, st.Message())
}

func TestFindByCpf_GetUserSuccessfully(t *testing.T) {
	req := &pb.UserByCpfRequest{
		Cpf: "440.072.470-05",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	objectID := objectIDFromHex("6404a984f18d899ec00c2a76")

	user := &model.User{
		ID:    objectID,
		Name:  "maria",
		Email: "maria@gmail.com",
		CPF:   req.GetCpf(),
	}

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	userRepoMock.On("FindByCpf", ctx, req.Cpf).Return(user, nil)

	service := service.UserService{UserRepository: userRepoMock}

	resp, err := service.FindByCpf(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), resp.GetId())
	assert.Equal(t, user.Email, resp.GetEmail())
	assert.Equal(t, user.Attributes, resp.GetAttributes())
}

func TestFindByEmail_InvalidEmail(t *testing.T) {
	req := &pb.UserByEmailRequest{
		Email: "invalid email",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.FindByEmail(ctx, req)

	st, _ := status.FromError(err)

	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestFindByEmail_UserNotFond(t *testing.T) {
	req := &pb.UserByEmailRequest{
		Email: "maria@gmail.com",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	userRepoMock.On("FindByEmail", ctx, req.Email).Return(nil, mongo.ErrNoDocuments)

	service := service.UserService{UserRepository: userRepoMock}

	_, err := service.FindByEmail(ctx, req)

	st, _ := status.FromError(err)

	msg := "user not found"

	assert.NotNil(t, err)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, msg, st.Message())
}

func TestFindByEmail_GetUserSuccessfully(t *testing.T) {
	req := &pb.UserByEmailRequest{
		Email: "maria@gmail.com",
	}

	ctx := context.Background()

	SetUpLog(ctx)

	objectID := objectIDFromHex("6404a984f18d899ec00c2a76")

	user := &model.User{
		ID:    objectID,
		Name:  "maria",
		Email: req.GetEmail(),
		CPF:   "440.072.470-05",
	}

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	userRepoMock.On("FindByEmail", ctx, req.GetEmail()).Return(user, nil)

	service := service.UserService{UserRepository: userRepoMock}

	resp, err := service.FindByEmail(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, user.ID.Hex(), resp.GetId())
	assert.Equal(t, user.Email, resp.GetEmail())
	assert.Equal(t, user.Attributes, resp.GetAttributes())
}

func SetUpLog(ctx context.Context) {
	cfg := SetUpConfig()
	log := logger.NewLog(cfg).GetGRPCLogger()
	log.WithContext(ctx)
}

func SetUpConfig() *config.Config {
	err := setEnvValues()
	if err != nil {
		panic(err)
	}
	var cfg config.Config
	cfg.MongoDB.URL = "localhost:20071"
	cfg.MongoDB.MongoDatabase = "test"
	cfg.MongoCollections.User.Collection = "test-user"

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	config.ExportConfig(&cfg)

	return &cfg
}

func setEnvValues() error {
	err := os.Setenv("APP_NAME", "user-manger-service")
	if err != nil {
		return fmt.Errorf("Error setting APP_NAME, err = %w", err)
	}

	err = os.Setenv("APP_VERSION", "1.0.0")
	if err != nil {
		return fmt.Errorf("Error setting APP_VERSION, err = %w", err)
	}

	err = os.Setenv("LOG_LEVEL", "debug")
	if err != nil {
		return fmt.Errorf("Error setting LOG_LEVEL, err = %w", err)
	}

	err = os.Setenv("GRPC_PORT", "50051")
	if err != nil {
		return fmt.Errorf("Error setting GRPC_PORT, err = %w", err)
	}

	return nil
}

func objectIDFromHex(id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	return objectID
}
