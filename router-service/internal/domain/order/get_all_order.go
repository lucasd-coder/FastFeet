package order

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/router-service/internal/shared/errors"
	"github.com/lucasd-coder/router-service/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceImpl) GetAllOrder(ctx context.Context, pld *GetAllOrderPayload) (*pb.GetAllOrderResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": pld,
	}).Info("received request")

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
		Addresses:     s.newGetAddress(pld),
	}

	res, err := s.businessRepo.GetAllOrder(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("fail call businessRepository err: %w", err)
	}

	return res, nil
}

func (s *ServiceImpl) newGetAddress(pld *GetAllOrderPayload) *pb.Address {
	return &pb.Address{
		Address:      pld.Address.Address,
		Number:       pld.Address.Number,
		PostalCode:   pld.Address.PostalCode,
		Neighborhood: pld.Address.Neighborhood,
		City:         pld.Address.City,
		State:        pld.Address.State,
	}
}
