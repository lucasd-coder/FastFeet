package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/order"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
)

type OrderHandler struct {
	authRepository      shared.AuthRepository
	viaCepRepository    model.ViaCepRepository
	orderDataRepository model.OrderDataRepository
	cfg                 *config.Config
	validate            shared.Validator
}

func NewOrderHandler(
	authRepo shared.AuthRepository,
	viaCepRepo model.ViaCepRepository,
	orderDatRepo model.OrderDataRepository,
	cfg *config.Config,
	validate shared.Validator) *OrderHandler {
	return &OrderHandler{
		authRepository:      authRepo,
		viaCepRepository:    viaCepRepo,
		orderDataRepository: orderDatRepo,
		cfg:                 cfg,
		validate:            validate,
	}
}

func (h *OrderHandler) Handler(ctx context.Context, m []byte) error {
	var pld model.Payload

	if err := json.Unmarshal(m, &pld); err != nil {
		return fmt.Errorf("err Unmarshal: %w", err)
	}
	return h.handler(ctx, pld)
}

func (h *OrderHandler) handler(ctx context.Context, pld model.Payload) error {
	log := logger.FromContext(ctx)

	fields := map[string]interface{}{
		"payload": pld,
	}

	log.WithFields(fields).Info("received payload")

	if err := pld.Validate(h.validate); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	err := h.hasActiveUser(ctx, pld.Data.DeliverymanID)
	if err != nil {
		return err
	}

	isAdmin, err := h.hasPermissionIsAdmin(ctx, pld.Data.UserID)
	if err != nil {
		return err
	}

	if !isAdmin {
		log.Errorf("error mission not permission to id: %s", pld.Data.UserID)
		return nil
	}

	log.Infof("get started address with postalCode: %s", pld.Data.Address.PostalCode)

	address, err := h.viaCepRepository.GetAddress(ctx, pld.Data.Address.PostalCode)
	if err != nil {
		log.Errorf("error when get address with postalCode: %s err: %v", pld.Data.Address.PostalCode, err)
		return err
	}

	if address.GetPostalCode() == "" {
		log.Errorf("error validating address invalid to payload: %v", pld)
		return nil
	}

	req := h.newOrderRequest(pld, address)

	resp, err := h.orderDataRepository.Save(ctx, req)
	if err != nil {
		log.Errorf("error while call order-repository err: %v", err)
		return err
	}

	log.Infof("event processed successfully id: %s generated", resp.GetId())

	return nil
}

func (h *OrderHandler) hasPermissionIsAdmin(ctx context.Context, id string) (bool, error) {
	log := logger.FromContext(ctx)

	log.Infof("get started roles with id: %s", id)

	roles, err := h.authRepository.FindRolesByID(ctx, id)
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

func (h *OrderHandler) hasActiveUser(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)

	log.Infof("get started to check is active user with id: %s", id)

	isActiveUser, err := h.authRepository.IsActiveUser(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrUserNotFound) {
			log.Errorf("%v with id %s", err, id)
			return nil
		}
		return err
	}

	if !isActiveUser.Active {
		log.Errorf("deliveryman not active with id: %s", id)
		return nil
	}
	return nil
}

func (h *OrderHandler) newOrderRequest(pld model.Payload, address *shared.ViaCepAddressResponse) *pb.OrderRequest {
	return &pb.OrderRequest{
		DeliverymanId: pld.Data.DeliverymanID,
		Product:       &pb.Product{Name: pld.Data.Product.Name},
		Address: &pb.Address{
			Address:      address.Address,
			PostalCode:   address.PostalCode,
			Neighborhood: address.Neighborhood,
			City:         address.City,
			State:        address.State,
			Number:       pld.Data.Address.Number,
		},
	}
}
