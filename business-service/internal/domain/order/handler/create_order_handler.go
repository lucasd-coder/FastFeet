package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	model "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) CreateOrderHandler(ctx context.Context, m []byte) error {
	var pld model.Payload

	if err := json.Unmarshal(m, &pld); err != nil {
		return fmt.Errorf("err Unmarshal: %w", err)
	}

	logger.FromContext(ctx).
		With(slog.Any("payload", pld)).Info("received payload")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("eventDate", pld.EventDate))
	span.SetAttributes(attribute.String("userId", pld.Data.UserID))
	span.SetAttributes(attribute.String("deliverymanId", pld.Data.DeliverymanID))

	resp, err := h.service.CreateOrder(ctx, pld)
	if err != nil {
		return err
	}

	logger.FromContext(ctx).
		Infof("event processed successfully id: %s generated", resp.GetId())

	return nil
}
