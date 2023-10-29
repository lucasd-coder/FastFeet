// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: request/user.proto

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
	UserHandler_FindUserByEmail_FullMethodName = "/pb.UserHandler/FindUserByEmail"
	UserHandler_GetRoles_FullMethodName        = "/pb.UserHandler/GetRoles"
	UserHandler_IsActiveUser_FullMethodName    = "/pb.UserHandler/IsActiveUser"
)

// UserHandlerClient is the client API for UserHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserHandlerClient interface {
	FindUserByEmail(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	GetRoles(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetRolesResponse, error)
	IsActiveUser(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*IsActiveUserResponse, error)
}

type userHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewUserHandlerClient(cc grpc.ClientConnInterface) UserHandlerClient {
	return &userHandlerClient{cc}
}

func (c *userHandlerClient) FindUserByEmail(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, UserHandler_FindUserByEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) GetRoles(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetRolesResponse, error) {
	out := new(GetRolesResponse)
	err := c.cc.Invoke(ctx, UserHandler_GetRoles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userHandlerClient) IsActiveUser(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*IsActiveUserResponse, error) {
	out := new(IsActiveUserResponse)
	err := c.cc.Invoke(ctx, UserHandler_IsActiveUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserHandlerServer is the server API for UserHandler service.
// All implementations must embed UnimplementedUserHandlerServer
// for forward compatibility
type UserHandlerServer interface {
	FindUserByEmail(context.Context, *EmptyRequest) (*GetUserResponse, error)
	GetRoles(context.Context, *EmptyRequest) (*GetRolesResponse, error)
	IsActiveUser(context.Context, *EmptyRequest) (*IsActiveUserResponse, error)
	mustEmbedUnimplementedUserHandlerServer()
}

// UnimplementedUserHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedUserHandlerServer struct {
}

func (UnimplementedUserHandlerServer) FindUserByEmail(context.Context, *EmptyRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindUserByEmail not implemented")
}
func (UnimplementedUserHandlerServer) GetRoles(context.Context, *EmptyRequest) (*GetRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoles not implemented")
}
func (UnimplementedUserHandlerServer) IsActiveUser(context.Context, *EmptyRequest) (*IsActiveUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsActiveUser not implemented")
}
func (UnimplementedUserHandlerServer) mustEmbedUnimplementedUserHandlerServer() {}

// UnsafeUserHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserHandlerServer will
// result in compilation errors.
type UnsafeUserHandlerServer interface {
	mustEmbedUnimplementedUserHandlerServer()
}

func RegisterUserHandlerServer(s grpc.ServiceRegistrar, srv UserHandlerServer) {
	s.RegisterService(&UserHandler_ServiceDesc, srv)
}

func _UserHandler_FindUserByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).FindUserByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserHandler_FindUserByEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).FindUserByEmail(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_GetRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).GetRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserHandler_GetRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).GetRoles(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserHandler_IsActiveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).IsActiveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserHandler_IsActiveUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).IsActiveUser(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserHandler_ServiceDesc is the grpc.ServiceDesc for UserHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserHandler",
	HandlerType: (*UserHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindUserByEmail",
			Handler:    _UserHandler_FindUserByEmail_Handler,
		},
		{
			MethodName: "GetRoles",
			Handler:    _UserHandler_GetRoles_Handler,
		},
		{
			MethodName: "IsActiveUser",
			Handler:    _UserHandler_IsActiveUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "request/user.proto",
}
