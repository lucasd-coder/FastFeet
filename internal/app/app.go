package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/provider/subscribe"
	"github.com/lucasd-coder/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/business-service/pkg/cache"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/monitor"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	logger := logger.NewLog(cfg)
	log := logger.GetGRPCLogger()

	cache.SetUpRedis(ctx, cfg)

	lis, err := net.Listen("tcp", "localhost:"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srvMetrics := monitor.RegisterSrvMetrics()
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	tp := monitor.RegisterOtel(ctx, cfg)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("Error shutting down tracer server provider: %v", err)
		}
	}()

	grpcServer := newGrpcServer(ctx, logger, reg)
	log.Infof("Started listening... address[:%s]", cfg.GRPC.Port)

	registerServices(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Could not serve: %v", err)
		}
	}()

	go newHTTPServer(ctx, cfg, reg)

	go subscribeUserEvents(ctx, cfg)

	go subscribeOrderEvents(ctx, cfg)

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	close(stopChan)

	grpcServer.GracefulStop()
}

func newGrpcServer(ctx context.Context, logger *logger.Log, reg prometheus.Registerer) *grpc.Server {
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	srvMetrics := monitor.RegisterSrvMetrics()

	grpcPanicRecoveryHandler := monitor.RegisterGrpcPanicRecoveryHandler(ctx, reg)

	interceptorOpt := otelgrpc.WithTracerProvider(otel.GetTracerProvider())

	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(interceptorOpt),
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logger.GetGRPCUnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(interceptorOpt),
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logger.GetGRPCStreamServerInterceptor(),
			grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)
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

	httpSrv := &http.Server{
		Addr:        ":" + cfg.HTTP.Port,
		ReadTimeout: cfg.HTTP.Timeout,
		Handler:     m,
	}
	log.Infof("starting HTTP server addr: %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Error(err)
		return
	}

	if err := httpSrv.Close(); err != nil {
		log.Error(err)
		return
	}
}

func subscribeUserEvents(ctx context.Context, cfg *config.Config) {
	optsQueueUserEvents := queueoptions.NewOptionQueueUserEvents(cfg)
	userHandler := InitializeUserHandler()
	subscribeUserEvents := subscribe.New(cfg.QueueUserEvents.URL, func(ctx context.Context, m []byte) error {
		if err := userHandler.Handler(ctx, m); err != nil {
			return err
		}
		return nil
	}, optsQueueUserEvents)

	subscribeUserEvents.Start(ctx)
}

func subscribeOrderEvents(ctx context.Context, cfg *config.Config) {
	optsQueueOrderEvents := queueoptions.NewOptionOrderEvents(cfg)
	orderHandler := InitializeOrderHandler()
	subscribeOrderEvents := subscribe.New(cfg.QueueOrderEvents.URL, func(ctx context.Context, m []byte) error {
		if err := orderHandler.Handler(ctx, m); err != nil {
			return err
		}
		return nil
	}, optsQueueOrderEvents)

	subscribeOrderEvents.Start(ctx)
}

func registerServices(grpcServer *grpc.Server) {
	orderDataService := InitializeOrderDataService()
	pb.RegisterGetAllOrderServiceServer(grpcServer, orderDataService)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)
}
