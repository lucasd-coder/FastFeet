package managerservice

import (
	"context"
	"time"

	"github.com/lucasd-coder/business-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func NewClient() (*grpc.ClientConn, error) {
	cfg := config.GetConfig()
	url := cfg.Integration.UserManagerService.URL

	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcRetryUnaryInterceptor),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func grpcRetryUnaryInterceptor(
	ctx context.Context,
	method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	cfg := config.GetConfig()
	maxRetry := cfg.Integration.UserManagerService.MaxRetry

	var err error

	for i := 0; i < maxRetry; i++ {
		err = invoker(ctx, method, req, reply, cc, opts...)
		if status.Code(err) != codes.Unavailable {
			break
		}
		time.Sleep(time.Second)
	}
	return err
}
