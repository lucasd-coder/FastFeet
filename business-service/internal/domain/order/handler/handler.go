package handler

import (
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/domain/order"
)

type Handler struct {
	service order.Service
	cfg     *config.Config
}

func NewHandler(s order.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
	}
}
