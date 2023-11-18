package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	appError "github.com/lucasd-coder/fast-feet/router-service/internal/shared/errors"
)

var (
	defaultLimit  int64 = 10
	defaultOffSet int64 = 1
	defaultNumber int64 = 0
)

type controller struct{}

func NewRouter(
	user *UserController,
	order *OrderController) *chi.Mux {
	r := chi.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", user.Save)
			r.Get("/{email}", user.FindUserByEmail)
		})
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
		log.Error("err during json.Marchal", err)
	}

	if _, err := w.Write(content); err != nil {
		log.Error("err during http.ResponseWriter", err)
	}
}

func (c *controller) getQueryParamConvertStringToInt(u *url.URL, param string, value int64) int64 {
	intValue, err := strconv.ParseInt(
		u.Query().Get(param), 10, 64)
	if err != nil {
		return value
	}
	return intValue
}
