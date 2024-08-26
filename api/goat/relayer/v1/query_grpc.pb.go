// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: goat/relayer/v1/query.proto

package relayerv1

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
	Query_Params_FullMethodName  = "/goat.relayer.v1.Query/Params"
	Query_Relayer_FullMethodName = "/goat.relayer.v1.Query/Relayer"
	Query_Voters_FullMethodName  = "/goat.relayer.v1.Query/Voters"
	Query_Pubkeys_FullMethodName = "/goat.relayer.v1.Query/Pubkeys"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Relayer queries relayer state
	Relayer(ctx context.Context, in *QueryRelayerRequest, opts ...grpc.CallOption) (*QueryRelayerResponse, error)
	// Voters queries all current voter list
	Voters(ctx context.Context, in *QueryVotersRequest, opts ...grpc.CallOption) (*QueryVotersResponse, error)
	// Pubkeys queries all current public keys
	Pubkeys(ctx context.Context, in *QueryPubkeysRequest, opts ...grpc.CallOption) (*QueryPubkeysResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Relayer(ctx context.Context, in *QueryRelayerRequest, opts ...grpc.CallOption) (*QueryRelayerResponse, error) {
	out := new(QueryRelayerResponse)
	err := c.cc.Invoke(ctx, Query_Relayer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Voters(ctx context.Context, in *QueryVotersRequest, opts ...grpc.CallOption) (*QueryVotersResponse, error) {
	out := new(QueryVotersResponse)
	err := c.cc.Invoke(ctx, Query_Voters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Pubkeys(ctx context.Context, in *QueryPubkeysRequest, opts ...grpc.CallOption) (*QueryPubkeysResponse, error) {
	out := new(QueryPubkeysResponse)
	err := c.cc.Invoke(ctx, Query_Pubkeys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Relayer queries relayer state
	Relayer(context.Context, *QueryRelayerRequest) (*QueryRelayerResponse, error)
	// Voters queries all current voter list
	Voters(context.Context, *QueryVotersRequest) (*QueryVotersResponse, error)
	// Pubkeys queries all current public keys
	Pubkeys(context.Context, *QueryPubkeysRequest) (*QueryPubkeysResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Relayer(context.Context, *QueryRelayerRequest) (*QueryRelayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Relayer not implemented")
}
func (UnimplementedQueryServer) Voters(context.Context, *QueryVotersRequest) (*QueryVotersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Voters not implemented")
}
func (UnimplementedQueryServer) Pubkeys(context.Context, *QueryPubkeysRequest) (*QueryPubkeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pubkeys not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Relayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRelayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Relayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Relayer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Relayer(ctx, req.(*QueryRelayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Voters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVotersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Voters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Voters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Voters(ctx, req.(*QueryVotersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pubkeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPubkeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pubkeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Pubkeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pubkeys(ctx, req.(*QueryPubkeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goat.relayer.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Relayer",
			Handler:    _Query_Relayer_Handler,
		},
		{
			MethodName: "Voters",
			Handler:    _Query_Voters_Handler,
		},
		{
			MethodName: "Pubkeys",
			Handler:    _Query_Pubkeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goat/relayer/v1/query.proto",
}