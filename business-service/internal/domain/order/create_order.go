package order

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (s *ServiceImpl) CreateOrder(ctx context.Context, pld Payload) (*pb.OrderResponse, error) {
	log := logger.FromContext(ctx)

	if err := pld.Validate(s.validate); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	err := s.hasActiveUser(ctx, pld.Data.DeliverymanID)
	if err != nil {
		return nil, err
	}

	isAdmin, err := s.hasPermissionIsAdmin(ctx, pld.Data.UserID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		log.Errorf("error mission not permission to id: %s", pld.Data.UserID)
		return nil, err
	}

	log.Infof("get started address with postalCode: %s", pld.Data.Address.PostalCode)

	address, err := s.viaCepRepository.GetAddress(ctx, pld.Data.Address.PostalCode)
	if err != nil {
		log.Errorf("error when get address with postalCode: %s err: %v", pld.Data.Address.PostalCode, err)
		return nil, err
	}

	if address.GetPostalCode() == "" {
		log.Error("error validating address invalid to", "payload", pld)
		return nil, err
	}

	req := s.newOrderRequest(pld, address)

	resp, err := s.orderRepository.Save(ctx, req)
	if err != nil {
		log.Errorf("error while call order-repository err: %v", err)
		return resp, err
	}

	return resp, nil
}

func (s *ServiceImpl) newOrderRequest(pld Payload, address *shared.ViaCepAddressResponse) *pb.OrderRequest {
	return &pb.OrderRequest{
		DeliverymanId: pld.Data.DeliverymanID,
		Product:       &pb.Product{Name: pld.Data.Product.Name},
		Addresses: &pb.Address{
			Address:      address.Address,
			PostalCode:   address.PostalCode,
			Neighborhood: address.Neighborhood,
			City:         address.City,
			State:        address.State,
			Number:       pld.Data.Address.Number,
		},
	}
}
