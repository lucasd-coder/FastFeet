package handler

import (
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
)

type Handler struct {
	service auth.Service
	cfg     *config.Config
}

func NewHandler(s auth.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
	}
}
