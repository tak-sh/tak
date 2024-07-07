// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: api/provider/v1beta1/provider.proto

package v1beta1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ProviderService_ListAccounts_FullMethodName         = "/tak.sh.api.provider.v1beta1.ProviderService/ListAccounts"
	ProviderService_Login_FullMethodName                = "/tak.sh.api.provider.v1beta1.ProviderService/Login"
	ProviderService_DownloadTransactions_FullMethodName = "/tak.sh.api.provider.v1beta1.ProviderService/DownloadTransactions"
)

// ProviderServiceClient is the client API for ProviderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderServiceClient interface {
	ListAccounts(ctx context.Context, in *ListAccounts_Request, opts ...grpc.CallOption) (*ListAccounts_Response, error)
	Login(ctx context.Context, in *Login_Request, opts ...grpc.CallOption) (*Login_Response, error)
	DownloadTransactions(ctx context.Context, in *DownloadTransactions_Request, opts ...grpc.CallOption) (*DownloadTransactions_Response, error)
}

type providerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProviderServiceClient(cc grpc.ClientConnInterface) ProviderServiceClient {
	return &providerServiceClient{cc}
}

func (c *providerServiceClient) ListAccounts(ctx context.Context, in *ListAccounts_Request, opts ...grpc.CallOption) (*ListAccounts_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAccounts_Response)
	err := c.cc.Invoke(ctx, ProviderService_ListAccounts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerServiceClient) Login(ctx context.Context, in *Login_Request, opts ...grpc.CallOption) (*Login_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Login_Response)
	err := c.cc.Invoke(ctx, ProviderService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerServiceClient) DownloadTransactions(ctx context.Context, in *DownloadTransactions_Request, opts ...grpc.CallOption) (*DownloadTransactions_Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownloadTransactions_Response)
	err := c.cc.Invoke(ctx, ProviderService_DownloadTransactions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderServiceServer is the server API for ProviderService service.
// All implementations should embed UnimplementedProviderServiceServer
// for forward compatibility
type ProviderServiceServer interface {
	ListAccounts(context.Context, *ListAccounts_Request) (*ListAccounts_Response, error)
	Login(context.Context, *Login_Request) (*Login_Response, error)
	DownloadTransactions(context.Context, *DownloadTransactions_Request) (*DownloadTransactions_Response, error)
}

// UnimplementedProviderServiceServer should be embedded to have forward compatible implementations.
type UnimplementedProviderServiceServer struct {
}

func (UnimplementedProviderServiceServer) ListAccounts(context.Context, *ListAccounts_Request) (*ListAccounts_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccounts not implemented")
}
func (UnimplementedProviderServiceServer) Login(context.Context, *Login_Request) (*Login_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedProviderServiceServer) DownloadTransactions(context.Context, *DownloadTransactions_Request) (*DownloadTransactions_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadTransactions not implemented")
}

// UnsafeProviderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProviderServiceServer will
// result in compilation errors.
type UnsafeProviderServiceServer interface {
	mustEmbedUnimplementedProviderServiceServer()
}

func RegisterProviderServiceServer(s grpc.ServiceRegistrar, srv ProviderServiceServer) {
	s.RegisterService(&ProviderService_ServiceDesc, srv)
}

func _ProviderService_ListAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccounts_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServiceServer).ListAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProviderService_ListAccounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServiceServer).ListAccounts(ctx, req.(*ListAccounts_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Login_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProviderService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServiceServer).Login(ctx, req.(*Login_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderService_DownloadTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadTransactions_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServiceServer).DownloadTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProviderService_DownloadTransactions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServiceServer).DownloadTransactions(ctx, req.(*DownloadTransactions_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// ProviderService_ServiceDesc is the grpc.ServiceDesc for ProviderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProviderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tak.sh.api.provider.v1beta1.ProviderService",
	HandlerType: (*ProviderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAccounts",
			Handler:    _ProviderService_ListAccounts_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _ProviderService_Login_Handler,
		},
		{
			MethodName: "DownloadTransactions",
			Handler:    _ProviderService_DownloadTransactions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/provider/v1beta1/provider.proto",
}
