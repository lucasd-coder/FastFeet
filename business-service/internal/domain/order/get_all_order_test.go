package order_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/business-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetAllOrderSuite struct {
	suite.Suite
	svc            order.Service
	repoAuth       *mocks.AuthRepository_internal_shared
	repoOrder      *mocks.Repository_internal_domain_order
	repoViaCep     *mocks.ViaCepRepository_internal_domain_order
	ctx            context.Context
	getAllOrderReq order.GetAllOrderRequest
}

func (suite *GetAllOrderSuite) SetupTest() {
	val := validator.NewValidation()
	repoAuth := new(mocks.AuthRepository_internal_shared)
	repoOrder := new(mocks.Repository_internal_domain_order)
	repoViaCep := new(mocks.ViaCepRepository_internal_domain_order)

	suite.repoAuth = repoAuth
	suite.repoOrder = repoOrder
	suite.repoViaCep = repoViaCep
	suite.svc = order.NewService(val, repoOrder, repoAuth, repoViaCep)
	suite.ctx = context.Background()
	suite.getAllOrderReq = order.GetAllOrderRequest{
		ID:            "656c916c3aa4eccdfb732a80",
		UserID:        "970ea619-4bc5-4d7a-9cfb-f5a775dde6f3",
		DeliverymanID: "004ae0f0-e4fa-44bf-8311-0030776205e7",
	}
}

func (suite *GetAllOrderSuite) TestGetAllOrderValidateFailure() {
	pld := order.GetAllOrderRequest{}

	_, err := suite.svc.GetAllOrder(suite.ctx, &pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.InvalidArgument)
}

func (suite *GetAllOrderSuite) TestGetAllOrderWhenUserInative() {
	pld := suite.getAllOrderReq
	respIsActiveUser := &shared.IsActiveUser{
		Active: false,
	}

	errUserUnauthorized := fmt.Errorf("%w: deliveryman not active with id: %s", shared.ErrUserUnauthorized, pld.DeliverymanID)

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.DeliverymanID).
		Return(respIsActiveUser, nil)

	_, err := suite.svc.GetAllOrder(suite.ctx, &pld)
	suite.Error(err)
	suite.Equal(err.Error(), errUserUnauthorized.Error())
}

func (suite *GetAllOrderSuite) TestGetAllOrderWhenUserRolesNotAdmin() {
	pld := suite.getAllOrderReq
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"USER"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.UserID).
		Return(getRolesResp, nil)

	_, err := suite.svc.GetAllOrder(suite.ctx, &pld)

	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.Unauthenticated)
}

func (suite *GetAllOrderSuite) TestGetAllOrderWhenUserNotFound() {
	pld := suite.getAllOrderReq

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.DeliverymanID).
		Return(nil, shared.NotFoundError(shared.ErrUserNotFound))

	_, err := suite.svc.GetAllOrder(suite.ctx, &pld)
	st, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(st.Code(), codes.NotFound)
}

func (suite *GetAllOrderSuite) TestGetAllOrder() {
	pld := suite.getAllOrderReq
	respIsActiveUser := &shared.IsActiveUser{
		Active: true,
	}

	suite.repoAuth.On("IsActiveUser", suite.ctx, pld.DeliverymanID).
		Return(respIsActiveUser, nil)

	getRolesResp := &shared.GetRolesResponse{
		Roles: []string{"ADMIN"},
	}

	suite.repoAuth.On("FindRolesByID", suite.ctx, pld.UserID).
		Return(getRolesResp, nil)

	getAllOrderResp := &pb.GetAllOrderResponse{
		Total:  1,
		Offset: 0,
		Limit:  10,
		Orders: []*pb.Order{
			{Id: "656c916c3aa4eccdfb732a80",
				DeliverymanId: pld.DeliverymanID,
				Product: &pb.Product{
					Name: "bola",
				}},
		},
	}

	suite.repoOrder.On("GetAllOrder", suite.ctx, mock.Anything).
		Return(getAllOrderResp, nil)

	getAll, err := suite.svc.GetAllOrder(suite.ctx, &pld)
	suite.NoError(err)
	suite.NotEmpty(getAll.GetOrders())
	suite.Equal(getAll.GetTotal(), int32(1))
}

func TestGetAllOrderSuite(t *testing.T) {
	suite.Run(t, new(GetAllOrderSuite))
}
