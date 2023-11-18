package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/internal/domain/user"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
)

type UserController struct {
	controller
	userService user.Service
}

func NewUserController(userService user.Service) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)

	pld := &user.User{}

	if err := json.NewDecoder(r.Body).Decode(pld); err != nil {
		msg := fmt.Errorf("error when doing decoder payload: %w", err)
		log.Error(msg.Error())
		h.SendError(ctx, w, msg)
		return
	}

	if err := h.userService.Save(ctx, pld); err != nil {
		h.SendError(ctx, w, err)
		return
	}

	resp := shared.CreateEvent{
		Message: "Please wait while we process your request.",
	}

	h.Response(ctx, w, resp, http.StatusOK)
}

func (h *UserController) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := chi.URLParam(r, "email")

	pld := user.FindByEmailRequest{
		Email: email,
	}

	resp, err := h.userService.FindUserByEmail(ctx, &pld)
	if err != nil {
		h.SendError(ctx, w, err)
		return
	}
	h.Response(ctx, w, resp, http.StatusOK)
}
