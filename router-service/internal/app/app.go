package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"log"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/bridge/opencensus"

	"github.com/lucasd-coder/fast-feet/router-service/config"
	"github.com/lucasd-coder/fast-feet/router-service/internal/controller"
	"github.com/lucasd-coder/fast-feet/router-service/internal/provider/middleware"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(cfg *config.Config) {
	optlogger := shared.NewOptLogger(cfg)
	optOtel := shared.NewOptOtel(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)

	tp, err := monitor.RegisterOtel(ctx, &optOtel)
	if err != nil {
		logDefault.Error("Error creating register otel ", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logDefault.Error("Error shutting down tracer server provider ", err)
		}
	}()
	opencensus.InstallTraceBridge()
	r := chi.NewRouter()
	r.Use(middleware.OpenTelemetryMiddleware(cfg.Name))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.PromMiddleware)

	logDefault.Info(fmt.Sprintf("Started listening... address[:%s]", cfg.Port))

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
		logDefault.Error(err.Error())
		return
	}
}
