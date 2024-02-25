package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserSuite struct {
	suite.Suite
	ctx      context.Context
	valErrs  noProviderVal.ValidationErrors
	repoAuth *mocks.AuthRepository_internal_shared
	repoUser *mocks.Repository_internal_domain_user
	svc      user.Service
}

func (suite *CreateUserSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoUser := new(mocks.Repository_internal_domain_user)

	suite.repoAuth = repoAuth
	suite.repoUser = repoUser
	suite.svc = user.NewService(repoUser, repoAuth, val)
	suite.ctx = context.Background()
}

func (suite *CreateUserSuite) TestCreateUser_ValidateFailure() {
	pld := &user.Payload{}

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *CreateUserSuite) TestCreateUser_ValidadeUserWithEmail() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Attributes: pld.Data.Attributes,
	}

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(userResp, nil)

	errUserAlreadyExist := fmt.Errorf("error validating user with email: %w", shared.ErrUserAlreadyExist)

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.EqualError(err, errUserAlreadyExist.Error())
}

func (suite *CreateUserSuite) TestCreateUser_ValidadeUserWithCpf() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Attributes: pld.Data.Attributes,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: pld.Data.CPF,
	}

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(userResp, nil)

	errUserAlreadyExist := fmt.Errorf("error validating user with cpf: %w", shared.ErrUserAlreadyExist)

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.EqualError(err, errUserAlreadyExist.Error())
}

func (suite *CreateUserSuite) TestCreateUser_RegisterAndReturnWhenUserAlreadyExist() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: pld.Data.CPF,
	}

	getUserResp := &shared.GetUserResponse{
		ID:       "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Email:    pld.Data.Email,
		Username: pld.Data.Email,
		Enabled:  true,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Attributes: pld.Data.Attributes,
	}

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	suite.repoAuth.On("FindByEmail", suite.ctx, pld.Data.Email).Return(getUserResp, nil)

	suite.repoUser.On("Save", suite.ctx, mock.Anything).Return(userResp, nil)

	resp, err := suite.svc.Save(suite.ctx, pld)
	suite.NoError(err)
	suite.Equal(resp.GetEmail(), pld.Data.Email)
	suite.Equal(resp.GetName(), pld.Data.Name)
}

func (suite *CreateUserSuite) TestCreateUser() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: pld.Data.CPF,
	}

	userResp := &pb.UserResponse{
		Id:         "46c77402-ba50-4b48-9bd9-1c4f97e36565",
		Name:       pld.Data.Name,
		Email:      pld.Data.Email,
		Attributes: pld.Data.Attributes,
	}

	register := &shared.RegisterUserResponse{
		ID: "46c77402-ba50-4b48-9bd9-1c4f97e36565",
	}

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	suite.repoAuth.On("FindByEmail", suite.ctx, pld.Data.Email).Return(&shared.GetUserResponse{}, nil)

	suite.repoAuth.On("Register", suite.ctx, mock.Anything).Return(register, nil)

	suite.repoUser.On("Save", suite.ctx, mock.Anything).Return(userResp, nil)

	resp, err := suite.svc.Save(suite.ctx, pld)
	suite.NoError(err)
	suite.Equal(resp.GetEmail(), pld.Data.Email)
	suite.Equal(resp.GetName(), pld.Data.Name)
}

func (suite *CreateUserSuite) TestCreateUser_ValidadeUserWithEmailWhenUnknownFailure() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	errUnknown := errors.New("error unknown")
	errResp := status.Errorf(codes.Unknown, "error: %s", errUnknown)

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(nil, errResp)

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.EqualError(err, errResp.Error())
}

func (suite *CreateUserSuite) TestCreateUser_ValidadeUserWithCpfWhenUnknownFailure() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: pld.Data.CPF,
	}

	errUnknown := errors.New("error unknown")
	errResp := status.Errorf(codes.Unknown, "error: %s", errUnknown)

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(nil, errResp)

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.EqualError(err, errResp.Error())
}

func (suite *CreateUserSuite) TestCreateUser_RegisterAndReturnWhenUserAlreadyExistWhenUnknownFailure() {
	pld := &user.Payload{
		Data: user.Data{
			Name:      "maria",
			Email:     "maria@gmail.com",
			CPF:       "857.484.630-91",
			Password:  "12345678@*",
			Authority: "USER",
		},
		EventDate: time.Now().Format(time.RFC3339),
	}

	userByEmailRequest := &pb.UserByEmailRequest{
		Email: pld.Data.Email,
	}

	userByCpfRequest := &pb.UserByCpfRequest{
		Cpf: pld.Data.CPF,
	}

	suite.repoUser.On("FindByEmail", suite.ctx, userByEmailRequest).Return(&pb.UserResponse{}, nil)

	suite.repoUser.On("FindByCpf", suite.ctx, userByCpfRequest).Return(&pb.UserResponse{}, nil)

	errUnknown := errors.New("error unknown")
	errResp := status.Errorf(codes.Unknown, "error: %s", errUnknown)

	suite.repoAuth.On("FindByEmail", suite.ctx, pld.Data.Email).Return(nil, errResp)

	_, err := suite.svc.Save(suite.ctx, pld)
	suite.Error(err)
	suite.EqualError(err, errResp.Error())
}

func TestCreateUserSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}
