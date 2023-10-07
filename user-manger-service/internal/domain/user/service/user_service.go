package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"github.com/lucasd-coder/user-manger-service/internal/domain/user/repository"
	pkgErrors "github.com/lucasd-coder/user-manger-service/internal/errors"
	pb "github.com/lucasd-coder/user-manger-service/pkg/pb"
	"github.com/lucasd-coder/user-manger-service/pkg/val"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	UserRepository model.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepo,
	}
}

func (service *UserService) Save(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	log.WithFields(logrus.Fields{
		"payload": req,
	}).Info("received request")

	pld := model.User{
		UserID:     req.GetUserId(),
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		CPF:        req.GetCpf(),
		Attributes: req.GetAttributes(),
		CreatedAt:  time.Now(),
	}

	if err := pld.Validate(); err != nil {
		return nil, pkgErrors.ValidationErrors(err)
	}

	user, err := service.UserRepository.FindByUserID(ctx, pld.UserID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	}

	if user != nil {
		msg := fmt.Sprintf("already exist user with userID: %s", pld.UserID)
		return nil, pkgErrors.AlreadyExistsError(msg)
	}

	user, err = service.UserRepository.Save(ctx, &pld)
	if err != nil {
		return nil, err
	}

	return buildUserResponse(user), nil
}

func (service *UserService) FindByCpf(ctx context.Context, req *pb.UserByCpfRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	if !val.IsCPF(req.GetCpf()) {
		violations := []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "cpf",
				Description: "invalid object cpf",
			},
		}
		return nil, pkgErrors.InvalidArgumentError(violations)
	}

	user, err := service.UserRepository.FindByCpf(ctx, req.GetCpf())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, pkgErrors.NotFoundError("user not found")
		}

		log.Errorf("failed to find customer with CPF in database. Error: %+v", err)
		return nil, err
	}

	log.Info("request findByCpf finished....")
	return buildUserResponse(user), nil
}

func (service *UserService) FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error) {
	log := logger.FromContext(ctx)

	if !validateEmail(req.GetEmail()) {
		violations := []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "email",
				Description: "invalid object email",
			},
		}
		return nil, pkgErrors.InvalidArgumentError(violations)
	}

	user, err := service.UserRepository.FindByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, pkgErrors.NotFoundError("user not found")
		}
		log.Errorf("failed to find customer with Email in database. Error: %+v", err)
		return nil, err
	}

	log.Info("request findByEmail finished....")
	return buildUserResponse(user), nil
}

func buildUserResponse(user *model.User) *pb.UserResponse {
	if user == nil {
		return nil
	}

	return &pb.UserResponse{
		UserId:     user.UserID,
		Name:       user.Name,
		Email:      user.Email,
		Attributes: user.Attributes,
		CreatedAt:  user.GetCreatedAt(),
	}
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
