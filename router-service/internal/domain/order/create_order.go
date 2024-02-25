package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
)

func (s *ServiceImpl) Save(ctx context.Context, order *Order) error {
	log := logger.FromContext(ctx)

	if err := order.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg.Error())
		return msg
	}

	eventDate := s.getEventDate()

	pld := Payload{
		Data:      *order,
		EventDate: eventDate,
	}

	enc, err := json.Marshal(pld)
	if err != nil {
		return fmt.Errorf("fail json.Marshal err: %w", err)
	}

	msg := shared.Message{
		Body: enc,
		Metadata: map[string]string{
			"language":   "en",
			"importance": "high",
		},
	}

	if err := s.publish.Send(ctx, &msg); err != nil {
		msg := fmt.Errorf("error publishing payload in queue: %w", err)
		log.Error(msg.Error())
		return msg
	}

	slog.With("payload", pld).Info("payload successfully processed")

	return nil
}
