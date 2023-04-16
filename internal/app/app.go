package app

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/lucasd-coder/user-manger-service/config"
	"github.com/lucasd-coder/user-manger-service/internal/domain/user/service"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	"github.com/lucasd-coder/user-manger-service/pkg/mongodb"
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

	lis, err := net.Listen("tcp", "localhost:"+cfg.Port)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				logger.GetGRPCUnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				logger.GetGRPCStreamServerInterceptor(),
				grpc_recovery.StreamServerInterceptor(),
			),
		),
	)

	userRepository := InitializeUserRepository()

	pb.RegisterUserServiceServer(grpcServer, &service.UserService{UserRepository: userRepository})

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	reflection.Register(grpcServer)

	log.Infof("Started listening... address[:%s]", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
