package app

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/business-service/internal/provider/subscribe"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	logger := logger.NewLog(cfg)

	log := logger.GetGRPCLogger()

	lis, err := net.Listen("tcp", "localhost:"+cfg.Port)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			logger.GetGRPCUnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			logger.GetGRPCStreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	reflection.Register(grpcServer)

	userRepository := InitializeUserRepository()

	userHandler := handler.NewUserHandler(userRepository, cfg)

	subscribe := subscribe.New(func(ctx context.Context, m []byte) error {
		if err := userHandler.Handler(ctx, m); err != nil {
			return err
		}
		return nil
	}, cfg)

	log.Infof("Started listening... address[:%s]", cfg.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Could not serve: %v", err)
		}
	}()

	subscribe.Start(ctx)
}
