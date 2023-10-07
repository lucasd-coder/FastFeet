package handler

import (
	"context"

	"github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/pkg/pb"
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
	log := logger.FromContext(ctx)

	log.WithFields(map[string]interface{}{
		"payload": req,
	}).Info("received request")

	pld := user.FindByEmailRequest{
		Email: req.GetEmail(),
	}

	resp, err := h.service.FindByEmail(ctx, &pld)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
