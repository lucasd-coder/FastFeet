package user

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared/ciphers"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared/codec"
)

func (s *ServiceImpl) Save(ctx context.Context, user *User) error {
	log := logger.FromContext(ctx)

	if err := user.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return msg
	}

	eventDate := s.getEventDate()

	pld := Payload{
		Data:      *user,
		EventDate: eventDate,
	}

	codec := codec.New[Payload]()

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
