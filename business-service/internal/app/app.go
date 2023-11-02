package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	// revive
	_ "net/http/pprof"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	orderHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/order/handler"
	userHandler "github.com/lucasd-coder/fast-feet/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/subscribe"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/utils"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/cache"
	"github.com/lucasd-coder/fast-feet/business-service/pkg/pb"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/lucasd-coder/fast-feet/pkg/profiler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	optlogger := shared.NewOptLogger(cfg)
	logger := logger.NewLog(optlogger)
	log := logger.GetLogger()

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache.SetUpRedis(ctx, cfg)
	optOtel := shared.NewOptOtel(cfg)
	tp, err := monitor.RegisterOtel(ctx, &optOtel)
	if err != nil {
		log.Errorf("Error creating register otel: %v", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down tracer server provider: %v", err)
		}
	}()
	reg := prometheus.NewRegistry()

	grpcServer := newGrpcServer(ctx, logger, reg)
	registerServices(grpcServer)

	log.Infof("Started listening... address[:%s]", cfg.GRPC.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Could not serve: %v", err)
		}
	}()

	go newHTTPServer(ctx, cfg, reg)

	go subscribeUserEvents(ctx, cfg, reg)

	go subscribeOrderEvents(ctx, cfg, reg)

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

	interceptorOpt := otelgrpc.WithTracerProvider(otel.GetTracerProvider())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.GetGRPCUnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(interceptorOpt),
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			logger.GetGRPCStreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(interceptorOpt),
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
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
		log.Error(err)
		return
	}

	if err := httpSrv.Close(); err != nil {
		log.Error(err)
		return
	}
}

func subscribeUserEvents(ctx context.Context, cfg *config.Config, reg *prometheus.Registry) {
	optsQueueUserEvents := queueoptions.NewOptionQueueUserEvents(cfg)
	userHandler := InitializeUserHandler()
	log := logger.FromContext(ctx)

	metric, err := monitor.CreateMetrics(utils.ExtractQueueName(cfg.QueueUserEvents.QueueURL), reg)
	if err != nil {
		log.Errorf("CreateMetrics Error: %s", err)
		return
	}

	subscribeUserEvents := subscribe.New(func(ctx context.Context, m []byte) error {
		return userHandler.CreateUser(ctx, m)
	}, optsQueueUserEvents, metric)

	subscribeUserEvents.Start(ctx)
}

func subscribeOrderEvents(ctx context.Context, cfg *config.Config, reg *prometheus.Registry) {
	optsQueueOrderEvents := queueoptions.NewOptionOrderEvents(cfg)
	orderHandler := InitializeOrderHandler()

	log := logger.FromContext(ctx)

	metric, err := monitor.CreateMetrics(utils.ExtractQueueName(cfg.QueueOrderEvents.QueueURL), reg)
	if err != nil {
		log.Errorf("CreateMetrics Error: %s", err)
		return
	}

	subscribeOrderEvents := subscribe.New(func(ctx context.Context, m []byte) error {
		return orderHandler.CreateOrderHandler(ctx, m)
	}, optsQueueOrderEvents, metric)

	subscribeOrderEvents.Start(ctx)
}

func registerServices(grpcServer *grpc.Server) {
	initializeOrder := InitializeOrderHandler()
	initializeUser := InitializeUserHandler()

	order := orderHandler.NewOrderHandler(*initializeOrder)
	user := userHandler.NewUserHandler(*initializeUser)

	pb.RegisterOrderHandlerServer(grpcServer, order)
	pb.RegisterUserHandlerServer(grpcServer, user)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)
}
