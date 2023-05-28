package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	appError "github.com/lucasd-coder/router-service/internal/shared/errors"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

type controller struct{}

func NewRouter(
	user *UserController,
	order *OrderController) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/users", user.Save)
	})

	r.Group(func(r chi.Router) {
		r.Route("/orders", func(r chi.Router) {
			r.Post("/{userId}", order.Save)
			r.Get("/{userId}", order.GetAllOrder)
		})
	})

	return r
}

func (c *controller) SendError(ctx context.Context, w http.ResponseWriter, err error) {
	errResp := appError.BuildError(err)

	c.Response(ctx, w, errResp, errResp.StatusCode)
}

func (c *controller) Response(ctx context.Context, w http.ResponseWriter, body interface{}, statusCode int) {
	log := logger.FromContext(ctx)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	content, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		msg := fmt.Errorf("err during json.Marchal: %w", err)
		log.Error(msg)
	}

	if _, err := w.Write(content); err != nil {
		msg := fmt.Errorf("err during http.ResponseWriter: %w", err)
		log.Error(msg)
	}
}
