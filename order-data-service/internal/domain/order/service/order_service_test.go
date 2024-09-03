package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lucasd-coder/fast-feet/order-data-service/config"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order/service"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/order-data-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderServiceSuite struct {
	suite.Suite
	cfg      config.Config
	id       string
	objectId primitive.ObjectID
	svc      service.OrderService
	repo     *mocks.OrderRepository_internal_domain_order
	ctx      context.Context
	valErrs  noProviderVal.ValidationErrors
}

func (suite *OrderServiceSuite) SetupSuite() {
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
	suite.ctx = context.Background()
}

func (suite *OrderServiceSuite) SetupTest() {
	val := validator.NewValidation()

	suite.repo = new(mocks.OrderRepository_internal_domain_order)
	suite.id = uuid.NewString()
	suite.objectId = primitive.NewObjectID()

	suite.svc = *service.NewOrderService(val, suite.repo)
}

func (suite *OrderServiceSuite) TearDownTest() {
	suite.repo.AssertExpectations(suite.T())
	suite.repo = nil
}

func (suite *OrderServiceSuite) TestSave() {
	tests := []struct {
		name     string
		args     *pb.OrderRequest
		mocks    func()
		check    func(name string, resp, wantResp *pb.OrderResponse, err error)
		wantResp *pb.OrderResponse
	}{
		{
			name: "test validation field deliverymanID",
			args: &pb.OrderRequest{
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Mossoró",
					State:        "Rio grande do norte",
				},
			},
			check: func(name string, resp, wantResp *pb.OrderResponse, err error) {
				if err != nil {
					st, ok := status.FromError(err)
					suite.True(ok, "%s ,ok = %v, err = %v", name, ok, err)
					suite.Equal(st.Code(), codes.InvalidArgument, "%s ,code = %d, wantCode = %d, err = %v", name, codes.InvalidArgument, st.Code(), err)
				}
			},
			wantResp: nil,
		},
		{
			name: "testing the error when saving order",
			args: &pb.OrderRequest{
				DeliverymanId: suite.id,
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Number:       30,
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Natal",
					State:        "Rio grande do norte",
				},
			},
			mocks: func() {
				suite.repo.On("Save", mock.Anything, mock.Anything).Return(nil, errors.New("bad gateway")).Once()
			},
			check: func(name string, resp, wantResp *pb.OrderResponse, err error) {
				suite.Error(err)
				suite.Empty(resp.GetCreatedAt())
			},
			wantResp: nil,
		},
		{
			name: "test successfully created",
			args: &pb.OrderRequest{
				DeliverymanId: suite.id,
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Number:       30,
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Natal",
					State:        "Rio grande do norte",
				},
			},
			mocks: func() {
				product := order.NewProduct("PC")
				Addresses := order.Address{
					Number:       30,
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Natal",
					State:        "Rio grande do norte",
				}

				pld := order.CreateOrder{
					DeliverymanID: suite.id,
					Product:       product,
					Address:       Addresses,
				}
				newOrder := order.NewOrder(pld)
				newOrder.ID = suite.objectId
				suite.repo.On("Save", mock.Anything, mock.Anything).Return(newOrder, nil).Once()
			},
			check: func(name string, resp, wantResp *pb.OrderResponse, err error) {
				suite.NoError(err)
				suite.Equal(wantResp.GetId(), resp.GetId())
			},
			wantResp: &pb.OrderResponse{
				Id: suite.objectId.Hex(),
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.mocks != nil {
				tt.mocks()
			}
			resp, err := suite.svc.Save(suite.ctx, tt.args)
			if tt.check != nil {
				tt.check(tt.name, resp, tt.wantResp, err)
			}
		})
	}
}

func (suite *OrderServiceSuite) TestGetAllOrder() {
	tests := []struct {
		name     string
		args     *pb.GetAllOrderRequest
		mocks    func()
		check    func(name string, resp, wantResp *pb.GetAllOrderResponse, err error)
		wantResp *pb.GetAllOrderResponse
	}{
		{
			name: "test validation field deliverymanID",
			args: &pb.GetAllOrderRequest{
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Mossoró",
					State:        "Rio grande do norte",
				},
			},
			check: func(name string, resp, wantResp *pb.GetAllOrderResponse, err error) {
				if err != nil {
					st, ok := status.FromError(err)
					suite.True(ok, "%s ,ok = %v, err = %v", name, ok, err)
					suite.Equal(st.Code(), codes.InvalidArgument, "%s ,code = %d, wantCode = %d, err = %v", name, codes.InvalidArgument, st.Code(), err)
				}
			},
			wantResp: nil,
		},
		{
			name: "testing the error when getting orders",
			args: &pb.GetAllOrderRequest{
				DeliverymanId: suite.id,
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Mossoró",
					State:        "Rio grande do norte",
				},
			},
			mocks: func() {
				suite.repo.On("FindAll", mock.Anything, mock.Anything).Return(nil, errors.New("bad gateway")).Once()
			},
			check: func(name string, resp, wantResp *pb.GetAllOrderResponse, err error) {
				suite.Error(err)
				suite.Empty(resp)
			},
			wantResp: nil,
		},
		{
			name: "test successfully getting orders",
			args: &pb.GetAllOrderRequest{
				DeliverymanId: suite.id,
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Mossoró",
					State:        "Rio grande do norte",
				},
			},
			mocks: func() {
				product := order.NewProduct("PC")
				Addresses := order.Address{
					Number:       30,
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Natal",
					State:        "Rio grande do norte",
				}

				pld := order.CreateOrder{
					DeliverymanID: suite.id,
					Product:       product,
					Address:       Addresses,
				}
				newOrder := order.NewOrder(pld)
				newOrder.ID = suite.objectId
				orders := []order.Order{*newOrder}
				suite.repo.On("FindAll", mock.Anything, mock.Anything).Return(orders, nil).Once()
			},
		},
		{
			name: "test empty order search",
			args: &pb.GetAllOrderRequest{
				DeliverymanId: suite.id,
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Mossoró",
					State:        "Rio grande do norte",
				},
			},
			mocks: func() {
				orders := []order.Order{}
				suite.repo.On("FindAll", mock.Anything, mock.Anything).Return(orders, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.mocks != nil {
				tt.mocks()
			}
			resp, err := suite.svc.GetAllOrder(suite.ctx, tt.args)
			if tt.check != nil {
				tt.check(tt.name, resp, tt.wantResp, err)
			}
		})
	}
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}
