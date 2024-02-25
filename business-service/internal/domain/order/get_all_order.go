package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) GetAllOrder(ctx context.Context, pld *GetAllOrderRequest) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	if err := pld.Validate(s.validate); err != nil {
		return nil, shared.ValidationErrors(err)
	}

	if err := s.hasActiveUser(ctx, pld.DeliverymanID); err != nil {
		return nil, err
	}

	isAdmin, err := s.hasPermissionIsAdmin(ctx, pld.UserID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		if !strings.EqualFold(pld.DeliverymanID, pld.UserID) {
			log.Errorf("error mission not permission to id: %s", pld.DeliverymanID)
			return nil, shared.UnauthenticatedError(shared.ErrUserUnauthorized)
		}
	}

	reqData := s.newGetOrderServiceAllOrderRequest(pld)

	resp, err := s.orderRepository.GetAllOrder(ctx, reqData)
	if err != nil {
		return nil, fmt.Errorf("error when call order-data err: %w", err)
	}

	return resp, nil
}

func (s *ServiceImpl) newGetOrderServiceAllOrderRequest(pld *GetAllOrderRequest) *pb.GetOrderServiceAllOrderRequest {
	address := &pb.Address{
		Address:      pld.Address.Address,
		Number:       pld.Address.Number,
		PostalCode:   pld.Address.PostalCode,
		Neighborhood: pld.Address.Neighborhood,
		City:         pld.Address.City,
		State:        pld.Address.State,
	}

	return &pb.GetOrderServiceAllOrderRequest{
		Id:            pld.ID,
		StartDate:     pld.StartDate,
		EndDate:       pld.EndDate,
		Product:       &pb.Product{Name: pld.Product.Name},
		Addresses:     address,
		CreatedAt:     pld.CreatedAt,
		UpdatedAt:     pld.UpdatedAt,
		DeliverymanId: pld.DeliverymanID,
		CanceledAt:    pld.CanceledAt,
		Limit:         pld.Limit,
		Offset:        pld.Offset,
	}
}
