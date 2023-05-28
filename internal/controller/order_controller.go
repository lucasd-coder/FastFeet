package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	model "github.com/lucasd-coder/router-service/internal/domain/order"
	"github.com/lucasd-coder/router-service/internal/domain/order/service"
	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type OrderController struct {
	controller
	orderService model.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (h *OrderController) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &model.CreateOrder{}

	if err := json.NewDecoder(r.Body).Decode(pld); err != nil {
		msg := fmt.Errorf("error when doing decoder payload: %w", err)
		log.Error(msg)
		h.SendError(ctx, w, msg)
		return
	}

	userID := chi.URLParam(r, "userId")

	order := pld.NewOrder(userID)

	if err := h.orderService.Save(ctx, order); err != nil {
		h.SendError(ctx, w, err)
		return
	}

	resp := shared.CreateEvent{
		Message: "Please wait while we process your request.",
	}

	h.Response(ctx, w, resp, http.StatusOK)
}

func (h *OrderController) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &model.GetAllOrderRequest{}

	if err := json.NewDecoder(r.Body).Decode(pld); err != nil {
		msg := fmt.Errorf("error when doing decoder payload: %w", err)
		log.Error(msg)
		h.SendError(ctx, w, msg)
		return
	}

	userID := chi.URLParam(r, "userId")

	pldGetAllPayload := &model.GetAllOrderPayload{
		GetAllOrderRequest: *pld,
		UserID:             userID,
	}

	resp, err := h.orderService.GetAllOrders(ctx, pldGetAllPayload)
	if err != nil {
		h.SendError(ctx, w, err)
		return
	}

	h.Response(ctx, w, resp, http.StatusOK)
}
