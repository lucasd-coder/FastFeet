package handler

import (
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/user"
)

type Handler struct {
	service user.Service
	cfg     *config.Config
}

func NewHandler(s user.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
	}
}
