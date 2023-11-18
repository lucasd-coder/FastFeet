package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	// revive
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/mongodb"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/lucasd-coder/fast-feet/pkg/profiler"
	"github.com/lucasd-coder/fast-feet/user-manger-service/config"
	"github.com/lucasd-coder/fast-feet/user-manger-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/user-manger-service/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	optlogger := shared.NewOptLogger(cfg)
	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	optMongo := shared.NewOptMongoDB(cfg)
	mongodb.SetUpMongoDB(ctx, &optMongo)

	defer func() {
		if err := mongodb.CloseConnMongoDB(ctx); err != nil {
			logDefault.Error("Unable to disconnect: ", err)
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Panic(err)
	}

	optOtel := shared.NewOptOtel(cfg)
	tp, err := monitor.RegisterOtel(ctx, &optOtel)
	if err != nil {
		logDefault.Error("Error creating register otel: ", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logDefault.Error("Error shutting down tracer server provider: ", err)
		}
	}()
	reg := prometheus.NewRegistry()
	grpcServer := newGrpcServer(ctx, logger, reg)
	registerServices(grpcServer)

	slog.Info(fmt.Sprintf("Started listening... address[:%s]", cfg.GRPC.Port))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Could not serve: %v", err)
		}
	}()

	go newHTTPServer(ctx, cfg, reg)

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	grpcServer.GracefulStop()
}

func newGrpcServer(ctx context.Context, logger *logger.Log, reg *prometheus.Registry) *grpc.Server {
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	srvMetrics := monitor.RegisterSrvMetrics()
	reg.MustRegister(srvMetrics)
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))

	grpcPanicRecoveryHandler := monitor.RegisterGrpcPanicRecoveryHandler(ctx, reg)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logger.GetLogUnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logger.GetLogStreamServerInterceptor(),
			grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	srvMetrics.InitializeMetrics(grpcServer)
	return grpcServer
}

func newHTTPServer(ctx context.Context, cfg *config.Config, reg prometheus.Gatherer) {
	log := logger.FromContext(ctx)

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			Timeout:           cfg.HTTP.Timeout,
		},
	))

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	profiler.StartProfiling(m)

	httpSrv := &http.Server{
		Addr:        ":" + cfg.HTTP.Port,
		ReadTimeout: cfg.HTTP.Timeout,
		Handler:     m,
	}
	log.Infof("starting HTTP server addr: %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Error(err.Error())
		return
	}

	if err := httpSrv.Close(); err != nil {
		log.Error(err.Error())
		return
	}
}

func registerServices(grpcServer *grpc.Server) {
	userService := InitializeUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)
}
