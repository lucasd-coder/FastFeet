package service

import (
	"context"

	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"github.com/lucasd-coder/user-manger-service/internal/errors"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	pb "github.com/lucasd-coder/user-manger-service/pkg/pb"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (service *UserService) Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	pld := model.User{
		ID:         "12345",
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		CPF:        req.GetCpf(),
		Attributes: req.GetAttributes(),
	}

	if err := pld.Validate(); err != nil {
		return nil, errors.ValidationErrors(err)
	}

	log.WithFields(logrus.Fields{
		"payload": req,
	}).Info("received request")

	return &pb.UserResponse{
		Email: pld.Email,
		Name:  pld.Name,
	}, nil
}
