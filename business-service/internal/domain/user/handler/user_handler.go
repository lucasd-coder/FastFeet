package handler

import (
	"context"
	"log/slog"

	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
)

type UserHandler struct {
	pb.UnimplementedUserHandlerServer
	Handler
}

func NewUserHandler(h Handler) *UserHandler {
	return &UserHandler{
		Handler: h,
	}
}

func (h *UserHandler) FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error) {
	slog.With("payload", req).
		Info("received request")

	pld := user.FindByEmailRequest{
		Email: req.GetEmail(),
	}

	resp, err := h.service.FindByEmail(ctx, &pld)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
