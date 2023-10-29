package handler

import (
	"context"
	"errors"
	"io"

	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
)

type UserHandler struct {
	pb.UnimplementedRegisterHandlerServer
	pb.UnimplementedUserHandlerServer
	Handler
}

func NewUserHandler(h Handler) *UserHandler {
	return &UserHandler{
		Handler: h,
	}
}

func (h *UserHandler) CreateUser(srv pb.RegisterHandler_CreateUserServer) error {
	ctx := srv.Context()
	log := logger.FromContext(ctx)
	log.Info("received stream request")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if errors.Is(err, io.EOF) {
			log.Info("finish stream request")
			return nil
		}

		if err != nil {
			log.Errorf("receive error %v", err)
			continue
		}
		pld := user.Register{
			FirstName: req.GetFirstName(),
			LastName:  req.GetLastName(),
			Username:  req.GetUsername(),
			Password:  req.GetPassword(),
			Roles:     req.GetAuthority().String(),
		}

		resp, err := h.service.CreateUser(ctx, &pld)
		if err != nil {
			return err
		}
		if err := srv.Send(resp); err != nil {
			log.Errorf("send stream error %v", err)
		}
	}
}

func (h *UserHandler) FindUserByEmail(ctx context.Context, _ *pb.EmptyRequest) (*pb.GetUserResponse, error) {
	email, err := getHeader(ctx, "email")
	if err != nil {
		return nil, err
	}

	pld := user.FindUserByEmail{
		Email: email,
	}

	return h.service.FindUserByEmail(ctx, &pld)
}

func (h *UserHandler) GetRoles(ctx context.Context, _ *pb.EmptyRequest) (*pb.GetRolesResponse, error) {
	id, err := getHeader(ctx, "id")
	if err != nil {
		return nil, err
	}

	pld := user.GetUserID{
		ID: id,
	}

	return h.service.GetRoles(ctx, &pld)
}

func (h *UserHandler) IsActiveUser(ctx context.Context, _ *pb.EmptyRequest) (*pb.IsActiveUserResponse, error) {
	id, err := getHeader(ctx, "id")
	if err != nil {
		return nil, err
	}

	pld := user.GetUserID{
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
				Description: name + "header invalid",
			},
		})
	}

	var value string
	if values := md.Get(name); len(values) > 0 {
		value = values[0]
	}
	return value, nil
}
