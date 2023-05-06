package service

import (
	"context"

	model "github.com/lucasd-coder/order-data-service/internal/domain/delivery"
	pkgErrors "github.com/lucasd-coder/order-data-service/internal/errors"
	"github.com/lucasd-coder/order-data-service/internal/provider/validator"
	"github.com/lucasd-coder/order-data-service/internal/shared"

	"github.com/lucasd-coder/order-data-service/pkg/logger"
	"github.com/lucasd-coder/order-data-service/pkg/pb"
)

type DeliveryService struct {
	pb.UnimplementedDeliveryServiceServer
	validate shared.Validator
}

func NewDeliveryService(validate *validator.Validation) *DeliveryService {
	return &DeliveryService{validate: validate}
}

func (s *DeliveryService) Save(ctx context.Context, req *pb.DeliveryRequest) (*pb.DeliveryResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": req,
	}).Info("received request")

	pld := model.Delivery{
		DeliverymanID: req.GetDeliverymanId(),
		Product:       model.NewProduct(req.GetProduct().GetName()),
		StartDate:     req.GetStartDate().AsTime(),
		EndDate:       req.GetEndDate().AsTime(),
		Address:       s.newAddress(req),
	}

	if err := pld.Validate(s.validate); err != nil {
		return nil, pkgErrors.ValidationErrors(err)
	}

	return &pb.DeliveryResponse{}, nil
}

func (s *DeliveryService) newAddress(req *pb.DeliveryRequest) model.Address {
	return model.Address{
		Address:      req.GetAddress().GetAddress(),
		Number:       req.GetAddress().GetNumber(),
		PostalCode:   req.GetAddress().GetPostalCode(),
		Neighborhood: req.GetAddress().GetNeighborhood(),
		City:         req.GetAddress().GetCity(),
		State:        req.GetAddress().GetState(),
	}
}
