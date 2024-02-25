package service_test

import (
	"context"
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order/service"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/order-data-service/pkg/pb"
	"github.com/stretchr/testify/suite"
)

type OrderServiceSuite struct {
	suite.Suite
	svc     service.OrderService
	repo    *mocks.OrderRepository_internal_domain_order
	ctx     context.Context
	valErrs noProviderVal.ValidationErrors
}

func (suite *OrderServiceSuite) SetupTest() {
	val := validator.NewValidation()
	repo := new(mocks.OrderRepository_internal_domain_order)

	suite.repo = repo
	suite.svc = *service.NewOrderService(val, repo)
	suite.ctx = context.Background()
}

func (suite *OrderServiceSuite) TestSaveValidation() {
	tests := []struct {
		name    string
		args    *pb.OrderRequest
		wantErr bool
	}{
		{
			name: "test validation field deliverymanID",
			args: &pb.OrderRequest{
				DeliverymanId: "1234567",
				Product: &pb.Product{
					Name: "mesa",
				},
				Addresses: &pb.Address{
					Address:      "rua das marias",
					PostalCode:   "123456",
					Neighborhood: "apt 10",
					City:         "Rio grande do norte",
					State:        "Rio grande do norte",
				},
			},
			wantErr: true,
		},
		{
			
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			_, err := suite.svc.Save(suite.ctx, tt.args)
			if err != nil {
				st, ok := status.FromError(err)
				suite.True(ok, "suite.svc.Save() = %v, wantErr %v", err, tt.wantErr)
				suite.Equal(st.Code(), codes.InvalidArgument, "suite.svc.Save() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}
