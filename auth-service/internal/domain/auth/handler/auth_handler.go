package handler

import (
	"context"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
)

type AuthHandler struct {
	pb.UnimplementedRegisterHandlerServer
	pb.UnimplementedAuthHandlerServer
	Handler
}

func NewAuthHandler(h Handler) *AuthHandler {
	return &AuthHandler{
		Handler: h,
	}
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log := logger.FromContext(ctx)
	log.Info("received request request")

	pld := auth.Register{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		Roles:     req.GetAuthority().String(),
	}

	return h.service.CreateUser(ctx, &pld)
}

func (h *AuthHandler) FindUserByEmail(ctx context.Context, _ *pb.EmptyRequest) (*pb.GetUserResponse, error) {
	email, err := getHeader(ctx, "email")
	if err != nil {
		return nil, err
	}

	pld := auth.FindUserByEmail{
		Email: email,
	}

	return h.service.FindUserByEmail(ctx, &pld)
}

func (h *AuthHandler) GetRoles(ctx context.Context, _ *pb.EmptyRequest) (*pb.GetRolesResponse, error) {
	id, err := getHeader(ctx, "id")
	if err != nil {
		return nil, err
	}

	pld := auth.GetUserID{
		ID: id,
	}

	return h.service.GetRoles(ctx, &pld)
}

func (h *AuthHandler) IsActiveUser(ctx context.Context, _ *pb.EmptyRequest) (*pb.IsActiveUserResponse, error) {
	id, err := getHeader(ctx, "id")
	if err != nil {
		return nil, err
	}

	pld := auth.GetUserID{
		ID: id,
	}

	return h.service.IsActiveUser(ctx, &pld)
}

func getHeader(ctx context.Context, name string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", shared.InvalidArgumentError([]*errdetails.BadRequest_FieldViolation{
			{
				Field:       name,
				Description: name + " header invalid",
			},
		})
	}

	var value string
	if values := md.Get(name); len(values) > 0 {
		value = values[0]
	}
	return value, nil
}
