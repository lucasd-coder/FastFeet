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
	"github.com/lucasd-coder/router-service/internal/shared/errors"
	"github.com/lucasd-coder/router-service/pkg/logger"
	"github.com/lucasd-coder/router-service/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	validate     shared.Validator
	publish      shared.Publish
	cfg          *config.Config
	businessRepo model.BusinessRepository
}

func NewOrderService(
	validate *validator.Validation,
	publish *publish.Published,
	cfg *config.Config,
	businessRepo model.BusinessRepository) *OrderService {
	return &OrderService{validate: validate, publish: publish, cfg: cfg, businessRepo: businessRepo}
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

func (s *OrderService) GetAllOrders(ctx context.Context, pld *model.GetAllOrderPayload) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	if err := pld.Validate(s.validate); err != nil {
		msg := fmt.Errorf("err validating payload: %w", err)
		log.Error(msg)
		return nil, msg
	}

	req := &pb.GetAllOrderRequest{
		Id:            pld.ID,
		UserId:        pld.UserID,
		DeliverymanId: pld.DeliverymanID,
		StartDate:     pld.StartDate,
		EndDate:       pld.EndDate,
		CreatedAt:     pld.CreatedAt,
		UpdatedAt:     pld.UpdatedAt,
		CanceledAt:    pld.CanceledAt,
		Limit:         pld.Limit,
		Offset:        pld.Offset,
		Product:       &pb.Product{Name: pld.Product.Name},
		Address:       s.newGetAddress(pld),
	}

	res, err := s.businessRepo.GetAllOrders(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("fail call businessRepository err: %w", err)
	}

	return res, nil
}

func (s *OrderService) newGetAddress(pld *model.GetAllOrderPayload) *pb.Address {
	return &pb.Address{
		Address:      pld.Address.Address,
		Number:       pld.Address.Number,
		PostalCode:   pld.Address.PostalCode,
		Neighborhood: pld.Address.Neighborhood,
		City:         pld.Address.City,
		State:        pld.Address.State,
	}
}
