package handler

import (
	"context"
	"log/slog"

	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
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
	logger.FromContext(ctx).With(slog.Any("payload", req)).Info("received request")

	pld := user.FindByEmailRequest{
		Email: req.GetEmail(),
	}

	resp, err := h.service.FindByEmail(ctx, &pld)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
