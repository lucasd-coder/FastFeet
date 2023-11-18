package handler

import (
	"context"
	"log/slog"

	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

type OrderHandler struct {
	pb.UnimplementedOrderHandlerServer
	Handler
}

func NewOrderHandler(h Handler) *OrderHandler {
	return &OrderHandler{
		Handler: h,
	}
}

func (g *OrderHandler) GetAllOrder(ctx context.Context, req *pb.GetAllOrderRequest) (
	*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	slog.With("payload", req).Info("received request")

	pld := g.newGetAllOrderRequest(req)

	resp, err := g.service.GetAllOrder(ctx, pld)
	if err != nil {
		return nil, err
	}

	log.Info("successfully fetching orders")

	return resp, nil
}

func (g *OrderHandler) newGetAllOrderRequest(req *pb.GetAllOrderRequest) *order.GetAllOrderRequest {
	address := order.GetAddress{
		Address:      req.GetAddresses().GetAddress(),
		Number:       req.GetAddresses().GetNumber(),
		PostalCode:   req.GetAddresses().GetPostalCode(),
		Neighborhood: req.GetAddresses().GetNeighborhood(),
		City:         req.GetAddresses().GetCity(),
		State:        req.GetAddresses().GetState(),
	}

	return &order.GetAllOrderRequest{
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
		Product:       order.GetProduct{Name: req.GetProduct().GetName()},
		Address:       address,
	}
}
