package handler

import (
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/domain/user"
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
