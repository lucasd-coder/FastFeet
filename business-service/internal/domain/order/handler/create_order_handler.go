package handler

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/lucasd-coder/business-service/internal/domain/order"
	"github.com/lucasd-coder/business-service/pkg/logger"
)

func (h *Handler) CreateOrderHandler(ctx context.Context, m []byte) error {
	log := logger.FromContext(ctx)

	var pld model.Payload

	if err := json.Unmarshal(m, &pld); err != nil {
		return fmt.Errorf("err Unmarshal: %w", err)
	}

	fields := map[string]interface{}{
		"payload": pld,
	}

	log.WithFields(fields).Info("received payload")

	resp, err := h.service.CreateOrder(ctx, pld)
	if err != nil {
		return err
	}

	log.Infof("event processed successfully id: %s generated", resp.GetId())

	return nil
}
