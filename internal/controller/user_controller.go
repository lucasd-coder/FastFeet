package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/lucasd-coder/router-service/internal/domain/user"
	"github.com/lucasd-coder/router-service/internal/domain/user/service"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type UserController struct {
	controller
	userService model.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(pld); err != nil {
		msg := fmt.Errorf("error when doing decoder payload: %w", err)
		log.Error(msg)
		h.SendError(ctx, w, msg)
		return
	}

	if err := h.userService.Save(ctx, pld); err != nil {
		h.SendError(ctx, w, err)
		return
	}

	h.Response(ctx, w, nil, http.StatusOK)
}
