package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

func (s *ServiceImpl) Save(ctx context.Context, order *Order) error {
	log := logger.FromContext(ctx)

	if err := order.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
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
		log.Error(msg)
		return msg
	}

	fields := map[string]interface{}{
		"payload": pld,
	}

	log.WithFields(fields).Info("payload successfully processed")

	return nil
}
