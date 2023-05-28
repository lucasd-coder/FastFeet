package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model "github.com/lucasd-coder/business-service/internal/domain/order"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
)

type OrderDataService struct {
	pb.UnimplementedGetAllOrderServiceServer
	validate        shared.Validator
	orderRepository model.OrderDataRepository
	authRepository  shared.AuthRepository
}

func NewOrderDataService(val shared.Validator,
	orderRepo model.OrderDataRepository,
	authRepo shared.AuthRepository) *OrderDataService {
	return &OrderDataService{
		validate:        val,
		orderRepository: orderRepo,
		authRepository:  authRepo,
	}
}

func (s *OrderDataService) GetAllOrders(ctx context.Context, req *pb.GetAllOrderRequest) (
	*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": req,
	}).Info("received request")

	pld := s.newGetAllOrderRequest(req)

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

	reqOrderService := s.newGetOrderServiceAllOrderRequest(pld)

	resp, err := s.orderRepository.GetAllOrders(ctx, reqOrderService)
	if err != nil {
		return nil, fmt.Errorf("error when call order-data err: %w", err)
	}

	return resp, nil
}

func (s *OrderDataService) newGetAllOrderRequest(req *pb.GetAllOrderRequest) *model.GetAllOrderRequest {
	address := model.GetAddress{
		Address:      req.GetAddresses().GetAddress(),
		Number:       req.GetAddresses().GetNumber(),
		PostalCode:   req.GetAddresses().GetPostalCode(),
		Neighborhood: req.GetAddresses().GetNeighborhood(),
		City:         req.GetAddresses().GetCity(),
		State:        req.GetAddresses().GetState(),
	}

	return &model.GetAllOrderRequest{
		ID:            req.GetId(),
		UserID:        req.GetUserId(),
		DeliverymanID: req.GetDeliverymanId(),
		StartDate:     req.GetStartDate(),
		EndDate:       req.GetEndDate(),
		CreatedAt:     req.GetCreatedAt(),
		UpdatedAt:     req.GetUpdatedAt(),
		CanceledAt:    req.GetCanceledAt(),
		Limit:         req.GetLimit(),
		Offset:        req.GetOffset(),
		Product:       model.GetProduct{Name: req.GetProduct().GetName()},
		Address:       address,
	}
}

func (s *OrderDataService) newGetOrderServiceAllOrderRequest(pld *model.GetAllOrderRequest) *pb.GetOrderServiceAllOrderRequest {
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

func (s *OrderDataService) hasPermissionIsAdmin(ctx context.Context, id string) (bool, error) {
	log := logger.FromContext(ctx)

	log.Infof("get started roles with id: %s", id)

	roles, err := s.authRepository.FindRolesByID(ctx, id)
	if err != nil {
		log.Errorf("error when check permission with id: %s, err: %v", id, err)
		return false, err
	}

	for _, role := range roles.Roles {
		if strings.EqualFold(shared.ADMIN, role) {
			return true, nil
		}
	}

	return false, nil
}

func (s *OrderDataService) hasActiveUser(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)

	log.Infof("get started to check is active user with id: %s", id)

	isActiveUser, err := s.authRepository.IsActiveUser(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrUserNotFound) {
			return shared.NotFoundError(shared.ErrUserNotFound)
		}
		return err
	}

	if !isActiveUser.Active {
		log.Errorf("deliveryman not active with id: %s", id)
		return fmt.Errorf("%w: deliveryman not active with id: %s", shared.ErrUserUnauthorized, id)
	}
	return nil
}
