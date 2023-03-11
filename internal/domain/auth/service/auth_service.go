package service

import (
	"context"

	"github.com/lucasd-coder/business-service/pkg/logger"
	pb "github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (service *AuthService) Save(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(logrus.Fields{
		"payload": req,
	}).Info("received request")

	return &pb.AuthResponse{
		Email:       "lucas@134@gmail.com",
		Deliveryman: true,
	}, nil
}
