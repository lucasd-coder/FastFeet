package handler_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/business-service/internal/mocks"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_InvalidPayload(t *testing.T) {
	model := &model.Payload{
		Name:       "",
		Email:      "test validate email",
		CPF:        "",
		Password:   "",
		Authority:  "",
		Attributes: map[string]string{},
	}

	pld, err := marshal(model)
	require.NoError(t, err)

	ctx := context.Background()

	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.NotNil(t, err)
}

func TestHandler_UserWithEmailAlreadyExist(t *testing.T) {
	model := &model.Payload{
		Name:       "maria",
		Email:      "maria@gmail.com",
		CPF:        "080.705.460-77",
		Password:   "123456",
		Authority:  "USER",
		Attributes: map[string]string{},
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: model.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       model.Name,
		Email:      model.Email,
		Attributes: model.Attributes,
	}

	pld, err := marshal(model)
	require.NoError(t, err)

	ctx := context.Background()

	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.ErrorIs(t, err, shared.ErrUserAlreadyExist)
}

func TestHandler_UserWithCPFAlreadyExist(t *testing.T) {
	model := &model.Payload{
		Name:       "maria",
		Email:      "maria@gmail.com",
		CPF:        "080.705.460-77",
		Password:   "123456",
		Authority:  "USER",
		Attributes: map[string]string{},
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: model.CPF,
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: model.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       model.Name,
		Email:      model.Email,
		Attributes: model.Attributes,
	}

	pld, err := marshal(model)
	require.NoError(t, err)

	ctx := context.Background()

	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	userRepoMock.On("FindByCpf", ctx, userByCpfRequest).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.ErrorIs(t, err, shared.ErrUserAlreadyExist)
}

func TestHandler_AuthAlreadyExist(t *testing.T) {
	payload := &model.Payload{
		Name:       "maria",
		Email:      "maria@gmail.com",
		CPF:        "080.705.460-77",
		Password:   "123456",
		Authority:  "USER",
		Attributes: map[string]string{},
	}
	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       payload.Name,
		Email:      payload.Email,
		Attributes: payload.Attributes,
	}
	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: payload.CPF,
	}
	userByEmailRequest := &pb.UserByEmailRequest{
		Email: payload.Email,
	}

	getUserResp := &model.GetUserResponse{
		ID:       "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Email:    payload.Email,
		Username: payload.Email,
		Enabled:  true,
	}

	pld, err := marshal(payload)
	require.NoError(t, err)

	ctx := context.Background()
	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	userRepoMock.On("FindByCpf", ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	authRepoMock.On("FindByEmail", ctx, payload.Email).Return(getUserResp, nil)

	userRepoMock.On("Save", ctx, mock.Anything).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.Nil(t, err)
}

func TestHandler_CreatedUserSuccessfully(t *testing.T) {
	payload := &model.Payload{
		Name:       "maria",
		Email:      "maria@gmail.com",
		CPF:        "080.705.460-77",
		Password:   "123456",
		Authority:  "USER",
		Attributes: map[string]string{},
	}
	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       payload.Name,
		Email:      payload.Email,
		Attributes: payload.Attributes,
	}
	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: payload.CPF,
	}
	userByEmailRequest := &pb.UserByEmailRequest{
		Email: payload.Email,
	}
	register := &model.RegisterUserResponse{
		ID: userResp.Id,
	}

	pld, err := marshal(payload)
	require.NoError(t, err)

	ctx := context.Background()
	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	userRepoMock.On("FindByCpf", ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	authRepoMock.On("FindByEmail", ctx, payload.Email).Return(&model.GetUserResponse{}, nil)

	authRepoMock.On("Register", ctx, payload.ToRegister()).Return(register, nil)

	userRepoMock.On("Save", ctx, mock.Anything).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.Nil(t, err)

}

func marshal(pld *model.Payload) ([]byte, error) {
	result, err := json.Marshal(pld)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SetUpLog(ctx context.Context) {
	cfg := SetUpConfig()
	log := logger.NewLog(cfg).GetGRPCLogger()
	log.WithContext(ctx)
}

func SetUpConfig() *config.Config {
	relativePath := "../../../../config/config-dev.yml"

	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	err = cleanenv.ReadConfig(absPath, &cfg)
	if err != nil {
		panic(err)
	}
	config.ExportConfig(&cfg)

	return &cfg
}
