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
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
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

type OrderHandlerSuite struct {
	suite.Suite
	srv          *grpc.Server
	lis          *bufconn.Listener
	cfg          config.Config
	orderHandler *handler.OrderHandler
	ctx          context.Context
	conn         *grpc.ClientConn
	repoAuth     *mocks.AuthRepository_internal_shared
	repoOrder    *mocks.Repository_internal_domain_order
	repoViaCep   *mocks.ViaCepRepository_internal_domain_order
}

func (suite *OrderHandlerSuite) SetupSuite() {
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

func (suite *OrderHandlerSuite) SetupTest() {
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
	repoOrder := new(mocks.Repository_internal_domain_order)
	repoViaCep := new(mocks.ViaCepRepository_internal_domain_order)

	suite.repoAuth = repoAuth
	suite.repoOrder = repoOrder
	suite.repoViaCep = repoViaCep
	svc := order.NewService(val, repoOrder, repoAuth, repoViaCep)
	hdler := handler.NewHandler(svc, &suite.cfg)
	suite.orderHandler = handler.NewOrderHandler(*hdler)

	pb.RegisterOrderHandlerServer(srv, suite.orderHandler)

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

func (suite *OrderHandlerSuite) TestGetAllOrder() {
	client := pb.NewOrderHandlerClient(suite.conn)

	deliverymanID := "1c42d3bf-6f10-40b6-94d6-e412c6287e5a"
	userID := "432280f4-2ed5-46ce-a0f1-c1984513dcdf"
	ID := "657c712d4cabfec758c31bbb"

	in := &pb.GetAllOrderRequest{
		UserId:        userID,
		Id:            ID,
		DeliverymanId: deliverymanID,
		Product: &pb.Product{
			Name: "bola",
		},
	}

	suite.repoAuth.On("IsActiveUser", mock.Anything, deliverymanID).Return(&shared.IsActiveUser{
		Active: true,
	}, nil)

	suite.repoAuth.On("FindRolesByID", mock.Anything, userID).Return(&shared.GetRolesResponse{
		Roles: []string{"admin"},
	}, nil)

	suite.repoOrder.On("GetAllOrder", mock.Anything, mock.Anything).Return(&pb.GetAllOrderResponse{
		Total:  1,
		Offset: 0,
		Limit:  10,
		Orders: []*pb.Order{
			{
				Id:            ID,
				DeliverymanId: deliverymanID,
				Product: &pb.Product{
					Name: "bola",
				},
			},
		},
	}, nil)

	resp, err := client.GetAllOrder(suite.ctx, in)
	suite.NoError(err)
	suite.Equal(resp.GetTotal(), int32(1))
	suite.NotEmpty(resp.GetOrders())
}

func (suite *OrderHandlerSuite) TestGetAllOrder_ValidateFailure() {
	client := pb.NewOrderHandlerClient(suite.conn)

	in := &pb.GetAllOrderRequest{}

	_, err := client.GetAllOrder(suite.ctx, in)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func TestOrderHandlerSuite(t *testing.T) {
	suite.Run(t, new(OrderHandlerSuite))
}
