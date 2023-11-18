package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared/ciphers"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared/codec"
)

func (s *ServiceImpl) Save(ctx context.Context, user *User) error {
	log := logger.FromContext(ctx)

	if err := user.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg.Error())
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
		log.Error(msg.Error())
		return msg
	}

	encrypt, err := ciphers.Encrypt(ciphers.ExtractKey([]byte(s.cfg.AesKey)), enc)
	if err != nil {
		msg := fmt.Errorf("err encrypting payload: %w", err)
		log.Error(msg.Error())
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
		log.Error(msg.Error())
		return msg
	}

	slog.With("payload",
		slog.String("name", pld.Data.Name),
		slog.String("eventDate", eventDate),
	).Info("payload successfully processed")

	return nil
}
