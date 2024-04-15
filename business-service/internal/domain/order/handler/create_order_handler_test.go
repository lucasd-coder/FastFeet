package handler_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateOrderHandlerSuite struct {
	suite.Suite
	cfg        config.Config
	ctx        context.Context
	handler    *handler.Handler
	repoAuth   *mocks.AuthRepository_internal_shared
	repoOrder  *mocks.Repository_internal_domain_order
	repoViaCep *mocks.ViaCepRepository_internal_domain_order
	valErrs    noProviderVal.ValidationErrors
}

func (suite *CreateOrderHandlerSuite) SetupSuite() {
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
}

func (suite *CreateOrderHandlerSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoOrder := new(mocks.Repository_internal_domain_order)
	repoViaCep := new(mocks.ViaCepRepository_internal_domain_order)

	suite.repoAuth = repoAuth
	suite.repoOrder = repoOrder
	suite.repoViaCep = repoViaCep

	svc := order.NewService(val, repoOrder, repoAuth, repoViaCep)
	suite.handler = handler.NewHandler(svc, &suite.cfg)
}

func (suite *CreateOrderHandlerSuite) TestCreateOrder_UnmarshalFailure() {
	body := []byte(`{""}`)

	err := suite.handler.CreateOrderHandler(suite.ctx, body)
	suite.Error(err)
}

func (suite *CreateOrderHandlerSuite) TestCreateOrder_ValidateFailure() {
	body := []byte(`{
	}`)

	err := suite.handler.CreateOrderHandler(suite.ctx, body)
	suite.Error(err)
	suite.ErrorAs(err, &suite.valErrs)
}

func (suite *CreateOrderHandlerSuite) TestCreateOrder() {
	body := []byte(`{
		"eventDate": "2023-11-05T17:59:26Z",
		"data": {
			"userId": "432280f4-2ed5-46ce-a0f1-c1984513dcdf",
			"deliverymanId": "7136f723-88dd-4f2f-8cf9-6207b65e7405",
			"product": {
				"name": "bola"
			},
			"address": {
				"postalCode": "01001000",
				"number": 10
			}
		}
	}`)

	deliverymanID := "7136f723-88dd-4f2f-8cf9-6207b65e7405"
	userID := "432280f4-2ed5-46ce-a0f1-c1984513dcdf"
	postalCode := "01001000"

	suite.repoAuth.On("IsActiveUser", suite.ctx, deliverymanID).Return(&shared.IsActiveUser{
		Active: true,
	}, nil)

	suite.repoAuth.On("FindRolesByID", suite.ctx, userID).Return(&shared.GetRolesResponse{
		Roles: []string{"admin"},
	}, nil)

	suite.repoViaCep.On("GetAddress", suite.ctx, postalCode).Return(&shared.AddressResponse{
		Address:      "rua das marias",
		PostalCode:   "01001000",
		Neighborhood: "parque dos Camargo",
		City:         "Jardim Europa",
		State:        "Rio grande do Sul",
	}, nil)

	suite.repoOrder.On("Save", suite.ctx, mock.Anything).Return(&pb.OrderResponse{
		Id:        "845343a4-0bd7-4918-94d2-fdbdb88c1679",
		CreatedAt: "2023-11-05T17:59:26Z",
	}, nil)

	err := suite.handler.CreateOrderHandler(suite.ctx, body)
	suite.NoError(err)
}

func TestCreateOrderHandlerSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderHandlerSuite))
}
