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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/lucasd-coder/user-manger-service/config"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	"github.com/lucasd-coder/user-manger-service/pkg/mongodb"
	"github.com/lucasd-coder/user-manger-service/pkg/monitor"
	"github.com/lucasd-coder/user-manger-service/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	logger := logger.NewLog(cfg)

	ctx := context.Background()

	log := logger.GetGRPCLogger()

	mongodb.SetUpMongoDB(ctx, cfg)

	defer func() {
		if err := mongodb.CloseConnMongoDB(ctx); err != nil {
			log.Errorf("Unable to disconnect: %v", err)
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", "localhost:"+cfg.GRPC.Port)
	if err != nil {
		panic(err)
	}

	srvMetrics := monitor.RegisterSrvMetrics()
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	tp, err := monitor.RegisterOtel(ctx, cfg)
	if err != nil {
		log.Errorf("Error creating register otel: %v", err)
		return
	}
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

	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

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

func registerServices(grpcServer *grpc.Server) {
	userService := InitializeUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)
}
