package handler_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/business-service/internal/mocks"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/internal/shared/ciphers"
	"github.com/lucasd-coder/business-service/internal/shared/codec"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_InvalidPayload(t *testing.T) {
	model := &model.Payload{
		Data: model.Data{
			Name:       "",
			Email:      "test validate email",
			CPF:        "",
			Password:   "",
			Authority:  "",
			Attributes: map[string]string{},
		},
	}

	pld, err := encode(model)
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
		Data: model.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now(),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: model.Data.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       model.Data.Name,
		Email:      model.Data.Email,
		Attributes: model.Data.Attributes,
	}

	pld, err := encode(model)
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
		Data: model.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now(),
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: model.Data.CPF,
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: model.Data.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       model.Data.Name,
		Email:      model.Data.Email,
		Attributes: model.Data.Attributes,
	}

	pld, err := encode(model)
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
		Data: model.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now(),
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

	getUserResp := &model.GetUserResponse{
		ID:       "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Email:    payload.Data.Email,
		Username: payload.Data.Email,
		Enabled:  true,
	}

	pld, err := encode(payload)
	require.NoError(t, err)

	ctx := context.Background()
	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	userRepoMock.On("FindByCpf", ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	authRepoMock.On("FindByEmail", ctx, payload.Data.Email).Return(getUserResp, nil)

	userRepoMock.On("Save", ctx, mock.Anything).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.Nil(t, err)
}

func TestHandler_CreatedUserSuccessfully(t *testing.T) {
	payload := &model.Payload{
		Data: model.Data{
			Name:       "maria",
			Email:      "maria@gmail.com",
			CPF:        "080.705.460-77",
			Password:   "123456",
			Authority:  "USER",
			Attributes: map[string]string{},
		},
		EventDate: time.Now(),
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
	register := &model.RegisterUserResponse{
		ID: userResp.Id,
	}

	pld, err := encode(payload)
	require.NoError(t, err)

	ctx := context.Background()
	SetUpLog(ctx)

	cfg := SetUpConfig()

	userRepoMock := new(mocks.UserRepository_internal_domain_user)
	authRepoMock := new(mocks.AuthRepository_internal_domain_user)

	userRepoMock.On("FindByEmail", ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	userRepoMock.On("FindByCpf", ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	authRepoMock.On("FindByEmail", ctx, payload.Data.Email).Return(&model.GetUserResponse{}, nil)

	authRepoMock.On("Register", ctx, payload.ToRegister()).Return(register, nil)

	userRepoMock.On("Save", ctx, mock.Anything).Return(userResp, nil)

	handler := handler.NewUserHandler(userRepoMock, authRepoMock, cfg)
	err = handler.Handler(ctx, pld)
	assert.Nil(t, err)
}

func encode(pld *model.Payload) ([]byte, error) {
	codec := codec.New[model.Payload]()

	enc, err := codec.Encode(*pld)
	if err != nil {
		return nil, err
	}

	cfg := SetUpConfig()
	result, err := ciphers.Encrypt(ciphers.ExtractKey([]byte(cfg.AesKey)), enc)
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
