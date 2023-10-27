package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
)

type OrderController struct {
	controller
	orderService order.Service
}

func NewOrderController(orderService order.Service) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (h *OrderController) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &order.CreateOrder{}

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

	pld := h.extractGetAllOrderRequest(r)

	userID := chi.URLParam(r, "userId")

	pldGetAllPayload := &order.GetAllOrderPayload{
		GetAllOrderRequest: *pld,
		UserID:             userID,
	}

	resp, err := h.orderService.GetAllOrder(ctx, pldGetAllPayload)
	if err != nil {
		log.Error(err)
		h.SendError(ctx, w, err)
		return
	}

	h.Response(ctx, w, resp, http.StatusOK)
}

func (h *OrderController) extractGetAllOrderRequest(r *http.Request) *order.GetAllOrderRequest {
	limit := h.getQueryParamConvertStringToInt(r.URL,
		"limit", defaultLimit)

	offset := h.getQueryParamConvertStringToInt(
		r.URL, "offset", defaultOffSet)

	number := h.getQueryParamConvertStringToInt(
		r.URL, "address.number", defaultNumber,
	)

	return &order.GetAllOrderRequest{
		ID:            r.URL.Query().Get("id"),
		DeliverymanID: r.URL.Query().Get("deliverymanId"),
		StartDate:     r.URL.Query().Get("startDate"),
		EndDate:       r.URL.Query().Get("endDate"),
		CreatedAt:     r.URL.Query().Get("createdAt"),
		UpdatedAt:     r.URL.Query().Get("updatedAt"),
		CanceledAt:    r.URL.Query().Get("canceledAt"),
		Limit:         limit,
		Offset:        offset,
		Product: order.GetProduct{
			Name: r.URL.Query().Get("product.name"),
		},
		Address: order.GetAddress{
			Address:      r.URL.Query().Get("address"),
			Number:       number,
			PostalCode:   r.URL.Query().Get("address.postalCode"),
			Neighborhood: r.URL.Query().Get("address.neighborhood"),
			City:         r.URL.Query().Get("address.city"),
			State:        r.URL.Query().Get("address.state"),
		},
	}
}
