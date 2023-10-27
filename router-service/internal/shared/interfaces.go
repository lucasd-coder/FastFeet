package shared

import (
	"context"

	"github.com/lucasd-coder/fast-feet/router-service/pkg/pb"
)

type (
	Validator interface {
		ValidateStruct(s interface{}) error
	}

	Publish interface {
		Send(ctx context.Context, msg *Message) error
	}

	BusinessRepository interface {
		GetAllOrder(ctx context.Context, req *pb.GetAllOrderRequest) (*pb.GetAllOrderResponse, error)
		FindByEmail(ctx context.Context, req *pb.UserByEmailRequest) (*pb.UserResponse, error)
	}
)
