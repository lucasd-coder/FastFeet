package app

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/internal/domain/user/handler"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
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
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			logger.GetGRPCStreamServerInterceptor(),
		),
	)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	reflection.Register(grpcServer)

	userRepository := InitializeUserRepository()

	userHandler := handler.NewUserHandler(userRepository)

	pld := &model.User{
		ID:    "640decf1cf68ef1c5cfc342d",
		Name:  "Maria",
		Email: "maria@gmail.com",
		CPF:   "872.169.470-29",
	}

	if err := userHandler.Create(context.TODO(), pld); err != nil {
		log.Fatal(err)
	}

	log.Infof("Started listening... address[:%s]", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}
}
