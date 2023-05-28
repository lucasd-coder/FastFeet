// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.9
// source: service/get_all_orders.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GetAllOrderService_GetAllOrders_FullMethodName = "/pb.GetAllOrderService/GetAllOrders"
)

// GetAllOrderServiceClient is the client API for GetAllOrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetAllOrderServiceClient interface {
	GetAllOrders(ctx context.Context, in *GetAllOrderRequest, opts ...grpc.CallOption) (*GetAllOrderResponse, error)
}

type getAllOrderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGetAllOrderServiceClient(cc grpc.ClientConnInterface) GetAllOrderServiceClient {
	return &getAllOrderServiceClient{cc}
}

func (c *getAllOrderServiceClient) GetAllOrders(ctx context.Context, in *GetAllOrderRequest, opts ...grpc.CallOption) (*GetAllOrderResponse, error) {
	out := new(GetAllOrderResponse)
	err := c.cc.Invoke(ctx, GetAllOrderService_GetAllOrders_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetAllOrderServiceServer is the server API for GetAllOrderService service.
// All implementations must embed UnimplementedGetAllOrderServiceServer
// for forward compatibility
type GetAllOrderServiceServer interface {
	GetAllOrders(context.Context, *GetAllOrderRequest) (*GetAllOrderResponse, error)
	mustEmbedUnimplementedGetAllOrderServiceServer()
}

// UnimplementedGetAllOrderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGetAllOrderServiceServer struct {
}

func (UnimplementedGetAllOrderServiceServer) GetAllOrders(context.Context, *GetAllOrderRequest) (*GetAllOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllOrders not implemented")
}
func (UnimplementedGetAllOrderServiceServer) mustEmbedUnimplementedGetAllOrderServiceServer() {}

// UnsafeGetAllOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetAllOrderServiceServer will
// result in compilation errors.
type UnsafeGetAllOrderServiceServer interface {
	mustEmbedUnimplementedGetAllOrderServiceServer()
}

func RegisterGetAllOrderServiceServer(s grpc.ServiceRegistrar, srv GetAllOrderServiceServer) {
	s.RegisterService(&GetAllOrderService_ServiceDesc, srv)
}

func _GetAllOrderService_GetAllOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetAllOrderServiceServer).GetAllOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetAllOrderService_GetAllOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetAllOrderServiceServer).GetAllOrders(ctx, req.(*GetAllOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GetAllOrderService_ServiceDesc is the grpc.ServiceDesc for GetAllOrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetAllOrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.GetAllOrderService",
	HandlerType: (*GetAllOrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllOrders",
			Handler:    _GetAllOrderService_GetAllOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/get_all_orders.proto",
}
