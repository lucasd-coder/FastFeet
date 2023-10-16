package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/controller"
	"github.com/lucasd-coder/router-service/internal/provider/middleware"
	"github.com/lucasd-coder/router-service/internal/shared"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(cfg *config.Config) {
	optlogger := shared.NewOptLogger(cfg)
	optOtel := shared.NewOptOtel(cfg)
	ctx := context.Background()

	logger := logger.NewLog(optlogger)
	log := logger.GetLogger()

	tp, err := monitor.RegisterOtel(ctx, &optOtel)
	if err != nil {
		log.Errorf("Error creating register otel: %v", err)
		return
	}
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
	r.Use(middleware.PromMiddleware)

	log.Infof("Started listening... address[:%s]", cfg.Port)

	userController := InitializeUserController()

	orderController := InitializeOrderController()

	controller := controller.NewRouter(userController, orderController)

	r.Mount("/", controller)
	r.Mount("/debug", chiMiddleware.Profiler())

	r.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			Registry:          prometheus.DefaultRegisterer,
			EnableOpenMetrics: true,
		},
	))

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
		return
	}

	if err := s.Close(); err != nil {
		log.Error(err)
		return
	}
}
