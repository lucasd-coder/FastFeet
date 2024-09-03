package order_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
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

type CreateOrderSuite struct {
	suite.Suite
	cfg        config.Config
	svc        order.Service
	repoAuth   *mocks.AuthRepository_internal_shared
	repoOrder  *mocks.Repository_internal_domain_order
	repoViaCep *mocks.ViaCepRepository_internal_domain_order
	ctx        context.Context
	valErrs    noProviderVal.ValidationErrors
	pld        order.Payload
}

func (suite *CreateOrderSuite) SetupSuite() {
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

func (suite *CreateOrderSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoOrder := new(mocks.Repository_internal_domain_order)
	repoViaCep := new(mocks.ViaCepRepository_internal_domain_order)

	suite.repoAuth = repoAuth
	suite.repoOrder = repoOrder
	suite.repoViaCep = repoViaCep
	suite.svc = order.NewService(val, repoOrder, repoAuth, repoViaCep)
	suite.ctx = context.Background()
	suite.pld = order.Payload{
		EventDate: time.Now().Format(time.RFC3339),
		Data: order.Data{
			UserID:        "970ea619-4bc5-4d7a-9cfb-f5a775dde6f3",
			DeliverymanID: "004ae0f0-e4fa-44bf-8311-0030776205e7",
			Product: order.Product{
				Name: "bola",
			},
			Address: order.Address{
				PostalCode: "12334567",
				Number:     20,
			},
		},
	}
}

func (suite *CreateOrderSuite) TestCreateOrderValidateFailure() {
	pld := order.Payload{}

	_, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.Error(err)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *CreateOrderSuite) TestCreateOrderWhenUserInative() {
	pld := suite.pld
	respIsActiveUser := &shared.IsActiveUser{
		Active: false,
	}

	errUserUnauthorized := fmt.Errorf("%w: deliveryman not active with id: %s", shared.ErrUserUnauthorized, pld.Data.DeliverymanID)

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(respIsActiveUser, nil)

	_, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.Error(err)
	suite.Equal(err.Error(), errUserUnauthorized.Error())
}

func (suite *CreateOrderSuite) TestCreateOrderWhenUserRolesNotAdmin() {
	pld := suite.pld
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"USER"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.Data.UserID).
		Return(getRolesResp, nil)

	resp, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.NoError(err)
	suite.Nil(resp)
}

func (suite *CreateOrderSuite) TestCreateOrderWhenUserNotfound() {
	pld := suite.pld

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(nil, shared.NotFoundError(shared.ErrUserNotFound))

	_, err := suite.svc.CreateOrder(suite.ctx, pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *CreateOrderSuite) TestCreateOrderWhenReturnViaCepRepositoryResponseError() {
	pld := suite.pld
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"ADMIN"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.Data.UserID).
		Return(getRolesResp, nil)

	suite.repoViaCep.On("GetAddress", suite.ctx,
		pld.Data.Address.PostalCode).Return(&shared.AddressResponse{}, shared.ErrExtractResponse)

	_, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.Error(err)
	suite.ErrorIs(err, shared.ErrExtractResponse)
}

func (suite *CreateOrderSuite) TestCreateOrderWhenAddressInvalid() {
	pld := suite.pld
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"ADMIN"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.Data.UserID).
		Return(getRolesResp, nil)

	suite.repoViaCep.On("GetAddress", suite.ctx, pld.Data.Address.PostalCode).Return(&shared.AddressResponse{}, nil)

	resp, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.NoError(err)
	suite.Nil(resp)
}

func (suite *CreateOrderSuite) TestCreateOrder() {
	pld := suite.pld
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.Data.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"ADMIN"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.Data.UserID).
		Return(getRolesResp, nil)

	getAddress := &shared.AddressResponse{
		Address:      "rua das marias",
		PostalCode:   "12345667",
		Neighborhood: "Copacabana",
		City:         "Rio de janeiro",
		State:        "Rio de janeiro",
	}

	suite.repoViaCep.On("GetAddress", suite.ctx,
		pld.Data.Address.PostalCode).Return(getAddress, nil)

	respOrderRepo := &pb.OrderResponse{
		Id:        "656caa24d0106f14d3aa2026",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	suite.repoOrder.On("Save", suite.ctx, mock.Anything).Return(respOrderRepo, nil)

	resp, err := suite.svc.CreateOrder(suite.ctx, pld)
	suite.NoError(err)
	suite.Equal(respOrderRepo, resp)
}

func (suite *CreateOrderSuite) TearDownTest() {
	suite.repoOrder.AssertExpectations(suite.T())
	suite.repoAuth.AssertExpectations(suite.T())
	suite.repoViaCep.AssertExpectations(suite.T())
}

func TestCreateOrderSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSuite))
}
