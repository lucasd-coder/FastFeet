package orderdataservice

import (
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"github.com/lucasd-coder/business-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg *config.Config) (*grpc.ClientConn, error) {
	url := cfg.Integration.OrderDataService.URL
	maxRetry := cfg.Integration.OrderDataService.MaxRetry

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(maxRetry),
		grpc_retry.WithPerRetryTimeout(1 * time.Second),
		grpc_retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded),
	}

	conn, err := grpc.Dial(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
