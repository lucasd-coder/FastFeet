package app

import (
	"context"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/provider/subscribe"
	"github.com/lucasd-coder/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/business-service/pkg/cache"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/lucasd-coder/business-service/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	cache.SetUpRedis(ctx, cfg)

	logger := logger.NewLog(cfg)

	log := logger.GetGRPCLogger()

	lis, err := net.Listen("tcp", "localhost:"+cfg.Port)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	userHandler := InitializeUserHandler()

	orderHandler := InitializeOrderHandler()

	orderDataService := InitializeOrderDataService()

	grpcServer := newGrpcServer(logger)

	pb.RegisterGetAllOrderServiceServer(grpcServer, orderDataService)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	reflection.Register(grpcServer)

	optsQueueUserEvents := queueoptions.NewOptionQueueUserEvents(cfg)

	optsQueueOrderEvents := queueoptions.NewOptionOrderEvents(cfg)

	subscribeUserEvents := subscribe.New(cfg.QueueUserEvents.URL, func(ctx context.Context, m []byte) error {
		if err := userHandler.Handler(ctx, m); err != nil {
			return err
		}
		return nil
	}, optsQueueUserEvents)

	subscribeOrderEvents := subscribe.New(cfg.QueueOrderEvents.URL, func(ctx context.Context, m []byte) error {
		if err := orderHandler.Handler(ctx, m); err != nil {
			return err
		}
		return nil
	}, optsQueueOrderEvents)

	log.Infof("Started listening... address[:%s]", cfg.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Could not serve: %v", err)
		}
	}()

	go subscribeUserEvents.Start(ctx)

	subscribeOrderEvents.Start(ctx)
}

func newGrpcServer(logger *logger.Log) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.GetGRPCUnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logger.GetGRPCStreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)
}
