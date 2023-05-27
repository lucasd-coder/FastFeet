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

	pld := &model.GetAllOrderRequest{
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
		Address:       s.newGetAddress(req),
	}

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

	reqOrderService := &pb.GetOrderServiceAllOrderRequest{
		Id:            req.GetId(),
		StartDate:     req.GetStartDate(),
		EndDate:       req.GetEndDate(),
		Product:       req.GetProduct(),
		Address:       req.GetAddress(),
		CreatedAt:     req.GetCreatedAt(),
		UpdatedAt:     req.GetUpdatedAt(),
		DeliverymanId: req.GetDeliverymanId(),
		CanceledAt:    req.GetCanceledAt(),
		Limit:         req.GetLimit(),
		Offset:        req.GetLimit(),
	}

	resp, err := s.orderRepository.GetAllOrders(ctx, reqOrderService)
	if err != nil {
		return nil, fmt.Errorf("error when call order-data err: %w", err)
	}

	return resp, nil
}

func (s *OrderDataService) newGetAddress(req *pb.GetAllOrderRequest) model.GetAddress {
	return model.GetAddress{
		Address:      req.GetAddress().GetAddress(),
		Number:       req.GetAddress().GetNumber(),
		PostalCode:   req.GetAddress().GetPostalCode(),
		Neighborhood: req.GetAddress().GetNeighborhood(),
		City:         req.GetAddress().GetCity(),
		State:        req.GetAddress().GetState(),
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
