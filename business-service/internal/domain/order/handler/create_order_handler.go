package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	model "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

func (h *Handler) CreateOrderHandler(ctx context.Context, m []byte) error {
	log := logger.FromContext(ctx)

	var pld model.Payload

	if err := json.Unmarshal(m, &pld); err != nil {
		return fmt.Errorf("err Unmarshal: %w", err)
	}

	slog.With("payload", pld).Info("received payload")

	resp, err := h.service.CreateOrder(ctx, pld)
	if err != nil {
		return err
	}

	log.Infof("event processed successfully id: %s generated", resp.GetId())

	return nil
}
