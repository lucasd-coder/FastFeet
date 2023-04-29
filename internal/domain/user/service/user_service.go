package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lucasd-coder/router-service/config"
	model "github.com/lucasd-coder/router-service/internal/domain/user"
	"github.com/lucasd-coder/router-service/internal/provider/publish"
	"github.com/lucasd-coder/router-service/internal/provider/validator"
	"github.com/lucasd-coder/router-service/internal/shared"
	"github.com/lucasd-coder/router-service/internal/shared/ciphers"
	"github.com/lucasd-coder/router-service/internal/shared/codec"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type UserService struct {
	validate shared.Validator
	publish  model.Publish
	cfg      *config.Config
}

func NewUserService(
	validate *validator.Validation,
	publish *publish.Published,
	cfg *config.Config) *UserService {
	return &UserService{validate: validate, publish: publish, cfg: cfg}
}

func (s *UserService) Save(ctx context.Context, user *model.User) error {
	log := logger.FromContext(ctx)

	if err := user.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return msg
	}

	eventDate := s.GetEventDate()

	pld := model.Payload{
		Data:      *user,
		EventDate: eventDate,
	}

	codec := codec.New[model.Payload]()

	enc, err := codec.Encode(pld)
	if err != nil {
		msg := fmt.Errorf("err encoding payload: %w", err)
		log.Error(msg)
		return msg
	}

	encrypt, err := ciphers.Encrypt(ciphers.ExtractKey([]byte(s.cfg.AesKey)), enc)
	if err != nil {
		msg := fmt.Errorf("err encrypting payload: %w", err)
		log.Error(msg)
		return msg
	}

	msg := shared.Message{
		Body: encrypt,
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
		"payload": map[string]string{
			"name":      pld.Data.Name,
			"eventDate": eventDate,
		},
	}

	log.WithFields(fields).Info("payload successfully processed")

	return nil
}

func (s *UserService) GetEventDate() string {
	return time.Now().Format(time.RFC3339)
}
