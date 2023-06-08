package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/controller"
	"github.com/lucasd-coder/router-service/internal/provider/middleware"
	"github.com/lucasd-coder/router-service/pkg/logger"
	"github.com/lucasd-coder/router-service/pkg/monitor"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	logger := logger.NewLog(cfg)

	log := logger.GetLogger()

	tp := monitor.RegisterOtel(ctx, cfg)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("Error shutting down tracer server provider: %v", err)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.OpenTelemetryMiddleware(cfg.Name))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)
	r.Use(chiMiddleware.Heartbeat("/" + cfg.Name + "/health"))

	log.Infof("Started listening... address[:%s]", cfg.Port)

	userController := InitializeUserController()

	orderController := InitializeOrderController()

	controller := controller.NewRouter(userController, orderController)

	r.Mount("/"+cfg.Name, controller)

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
