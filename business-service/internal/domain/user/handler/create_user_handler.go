package handler

import (
	"context"
	"fmt"
	"log/slog"

	model "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/ciphers"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) CreateUser(ctx context.Context, m []byte) error {
	var pld model.Payload
	log := logger.FromContext(ctx)

	dec, err := ciphers.Decrypt(ciphers.ExtractKey([]byte(h.cfg.AesKey)), m)
	if err != nil {
		return fmt.Errorf("err Decrypt: %w", err)
	}

	codec := codec.New[model.Payload]()

	if err := codec.Decode(dec, &pld); err != nil {
		return fmt.Errorf("err Decode: %w", err)
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("eventDate", pld.EventDate))
	span.SetAttributes(attribute.String("userName", pld.Data.Name))

	log.With(slog.Group("payload",
		"name", pld.Data.Name)).
		Info("received payload")

	user, err := h.service.Save(ctx, &pld)
	if err != nil {
		return err
	}

	log.Infof("payload successfully processed for id: %s", user.GetId())

	return nil
}
