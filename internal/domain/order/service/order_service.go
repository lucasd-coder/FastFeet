package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucasd-coder/router-service/config"
	model "github.com/lucasd-coder/router-service/internal/domain/order"
	"github.com/lucasd-coder/router-service/internal/provider/publish"
	"github.com/lucasd-coder/router-service/internal/provider/validator"
	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type OrderService struct {
	validate shared.Validator
	publish  shared.Publish
	cfg      *config.Config
}

func NewOrderService(
	validate *validator.Validation,
	publish *publish.Published,
	cfg *config.Config) *OrderService {
	return &OrderService{validate: validate, publish: publish, cfg: cfg}
}

func (s *OrderService) Save(ctx context.Context, order *model.Order) error {
	log := logger.FromContext(ctx)

	if err := order.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return msg
	}

	eventDate := s.GetEventDate()

	pld := model.Payload{
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

func (s *OrderService) GetEventDate() string {
	return time.Now().Format(time.RFC3339)
}
