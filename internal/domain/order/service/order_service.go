package service

import (
	"context"
	"fmt"

	model "github.com/lucasd-coder/order-data-service/internal/domain/order"
	"github.com/lucasd-coder/order-data-service/internal/domain/order/repository"
	pkgErrors "github.com/lucasd-coder/order-data-service/internal/errors"
	"github.com/lucasd-coder/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/order-data-service/internal/shared"

	"github.com/lucasd-coder/order-data-service/pkg/logger"
	"github.com/lucasd-coder/order-data-service/pkg/pb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	validate        shared.Validator
	orderRepository model.OrderRepository
}

func NewOrderService(
	validate *validator.Validation,
	orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{validate: validate, orderRepository: orderRepo}
}

func (s *OrderService) Save(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": req,
	}).Info("received request")

	pld := model.CreateOrder{
		DeliverymanID: req.GetDeliverymanId(),
		Product:       model.NewProduct(req.GetProduct().GetName()),
		Address:       s.newAddress(req),
	}

	if err := pld.Validate(s.validate); err != nil {
		return nil, pkgErrors.ValidationErrors(err)
	}

	order := model.NewOrder(pld)

	newOrder, err := s.orderRepository.Save(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("error when save: %w", err)
	}

	log.Infof("successfully created order with id: %s", newOrder.ID.Hex())

	return &pb.OrderResponse{
		Id:        newOrder.ID.Hex(),
		CreatedAt: order.GetCreatedAt(),
	}, nil
}

func (s *OrderService) newAddress(req *pb.OrderRequest) model.Address {
	return model.Address{
		Address:      req.GetAddresses().GetAddress(),
		Number:       req.GetAddresses().GetNumber(),
		PostalCode:   req.GetAddresses().GetPostalCode(),
		Neighborhood: req.GetAddresses().GetNeighborhood(),
		City:         req.GetAddresses().GetCity(),
		State:        req.GetAddresses().GetState(),
	}
}

func (s *OrderService) GetAllOrder(ctx context.Context, req *pb.GetAllOrderRequest) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": req,
	}).Info("received request")

	pld := &model.GetAllOrderRequest{
		ID:            req.GetId(),
		DeliverymanID: req.GetDeliverymanId(),
		StartDate:     req.GetStartDate(),
		EndDate:       req.GetEndDate(),
		CreatedAt:     req.GetCreatedAt(),
		UpdatedAt:     req.GetUpdatedAt(),
		CanceledAt:    req.GetCanceledAt(),
		Limit:         req.GetLimit(),
		Offset:        req.GetOffset(),
		Product:       model.GetProduct{Name: req.GetProduct().GetName()},
		Address:       s.newGetAddress(req),
	}

	if err := pld.Validate(s.validate); err != nil {
		return nil, pkgErrors.ValidationErrors(err)
	}

	orders, err := s.orderRepository.FindAll(ctx, pld)
	if err != nil {
		return nil, fmt.Errorf("error when orderRepository findAll: %w", err)
	}

	log.Info("successfully return getAllOrder")

	return s.extractGetAllOrderResponse(pld, orders), nil
}

func (s *OrderService) newGetAddress(req *pb.GetAllOrderRequest) model.GetAddress {
	return model.GetAddress{
		Address:      req.GetAddresses().GetAddress(),
		Number:       req.GetAddresses().GetNumber(),
		PostalCode:   req.GetAddresses().GetPostalCode(),
		Neighborhood: req.GetAddresses().GetNeighborhood(),
		City:         req.GetAddresses().GetCity(),
		State:        req.GetAddresses().GetState(),
	}
}

func (s *OrderService) extractGetAllOrderResponse(pld *model.GetAllOrderRequest, orders []model.Order) *pb.GetAllOrderResponse {
	if len(orders) == 0 {
		return &pb.GetAllOrderResponse{}
	}

	pbOrders := []*pb.Order{}

	for _, order := range orders {
		pbOrders = append(pbOrders, s.extractPbOrder(order))
	}

	return &pb.GetAllOrderResponse{
		Total:  int32(len(pbOrders)),
		Offset: int32(pld.Offset),
		Limit:  int32(pld.Limit),
		Orders: pbOrders,
	}
}

func (s *OrderService) extractPbOrder(order model.Order) *pb.Order {
	return &pb.Order{
		Id:            order.ID.Hex(),
		DeliverymanId: order.DeliverymanID,
		StartDate:     order.GetStartDate(),
		EndDate:       order.GetStartDate(),
		Product:       &pb.Product{Name: order.Product.Name},
		Addresses: &pb.Address{
			Address:      order.Address.Address,
			PostalCode:   order.Address.PostalCode,
			Neighborhood: order.Address.Neighborhood,
			City:         order.Address.City,
			State:        order.Address.State,
			Number:       order.Address.Number,
		},
		CreatedAt:  order.GetCreatedAt(),
		UpdatedAt:  order.GetUpdatedAt(),
		CanceledAt: order.GetCanceledAt(),
	}
}
