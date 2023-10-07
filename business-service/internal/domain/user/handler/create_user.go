package handler

import (
	"context"
	"fmt"

	model "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/ciphers"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
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

	fields := map[string]interface{}{
		"payload": map[string]string{
			"name": pld.Data.Name,
		},
	}
	log.WithFields(fields).Info("received payload")

	user, err := h.service.Save(ctx, &pld)
	if err != nil {
		return err
	}

	log.Infof("payload successfully processed for id: %s", user.Id)

	return nil
}
