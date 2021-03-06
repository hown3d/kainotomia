// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: kainotomia/v1alpha1/kainotomia.proto

package kainotomia

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

// KainotomiaServiceClient is the client API for KainotomiaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KainotomiaServiceClient interface {
	CreatePlaylist(ctx context.Context, in *CreatePlaylistRequest, opts ...grpc.CallOption) (*CreatePlaylistResponse, error)
	TriggerUpdate(ctx context.Context, in *TriggerUpdateRequest, opts ...grpc.CallOption) (*TriggerUpdateResponse, error)
	DeletePlaylist(ctx context.Context, in *DeletePlaylistRequest, opts ...grpc.CallOption) (*DeletePlaylistResponse, error)
}

type kainotomiaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKainotomiaServiceClient(cc grpc.ClientConnInterface) KainotomiaServiceClient {
	return &kainotomiaServiceClient{cc}
}

func (c *kainotomiaServiceClient) CreatePlaylist(ctx context.Context, in *CreatePlaylistRequest, opts ...grpc.CallOption) (*CreatePlaylistResponse, error) {
	out := new(CreatePlaylistResponse)
	err := c.cc.Invoke(ctx, "/kainotomia.v1alpha1.KainotomiaService/CreatePlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kainotomiaServiceClient) TriggerUpdate(ctx context.Context, in *TriggerUpdateRequest, opts ...grpc.CallOption) (*TriggerUpdateResponse, error) {
	out := new(TriggerUpdateResponse)
	err := c.cc.Invoke(ctx, "/kainotomia.v1alpha1.KainotomiaService/TriggerUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kainotomiaServiceClient) DeletePlaylist(ctx context.Context, in *DeletePlaylistRequest, opts ...grpc.CallOption) (*DeletePlaylistResponse, error) {
	out := new(DeletePlaylistResponse)
	err := c.cc.Invoke(ctx, "/kainotomia.v1alpha1.KainotomiaService/DeletePlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KainotomiaServiceServer is the server API for KainotomiaService service.
// All implementations should embed UnimplementedKainotomiaServiceServer
// for forward compatibility
type KainotomiaServiceServer interface {
	CreatePlaylist(context.Context, *CreatePlaylistRequest) (*CreatePlaylistResponse, error)
	TriggerUpdate(context.Context, *TriggerUpdateRequest) (*TriggerUpdateResponse, error)
	DeletePlaylist(context.Context, *DeletePlaylistRequest) (*DeletePlaylistResponse, error)
}

// UnimplementedKainotomiaServiceServer should be embedded to have forward compatible implementations.
type UnimplementedKainotomiaServiceServer struct {
}

func (UnimplementedKainotomiaServiceServer) CreatePlaylist(context.Context, *CreatePlaylistRequest) (*CreatePlaylistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePlaylist not implemented")
}
func (UnimplementedKainotomiaServiceServer) TriggerUpdate(context.Context, *TriggerUpdateRequest) (*TriggerUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerUpdate not implemented")
}
func (UnimplementedKainotomiaServiceServer) DeletePlaylist(context.Context, *DeletePlaylistRequest) (*DeletePlaylistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePlaylist not implemented")
}

// UnsafeKainotomiaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KainotomiaServiceServer will
// result in compilation errors.
type UnsafeKainotomiaServiceServer interface {
	mustEmbedUnimplementedKainotomiaServiceServer()
}

func RegisterKainotomiaServiceServer(s grpc.ServiceRegistrar, srv KainotomiaServiceServer) {
	s.RegisterService(&KainotomiaService_ServiceDesc, srv)
}

func _KainotomiaService_CreatePlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePlaylistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KainotomiaServiceServer).CreatePlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kainotomia.v1alpha1.KainotomiaService/CreatePlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KainotomiaServiceServer).CreatePlaylist(ctx, req.(*CreatePlaylistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KainotomiaService_TriggerUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KainotomiaServiceServer).TriggerUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kainotomia.v1alpha1.KainotomiaService/TriggerUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KainotomiaServiceServer).TriggerUpdate(ctx, req.(*TriggerUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KainotomiaService_DeletePlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePlaylistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KainotomiaServiceServer).DeletePlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kainotomia.v1alpha1.KainotomiaService/DeletePlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KainotomiaServiceServer).DeletePlaylist(ctx, req.(*DeletePlaylistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KainotomiaService_ServiceDesc is the grpc.ServiceDesc for KainotomiaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KainotomiaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kainotomia.v1alpha1.KainotomiaService",
	HandlerType: (*KainotomiaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePlaylist",
			Handler:    _KainotomiaService_CreatePlaylist_Handler,
		},
		{
			MethodName: "TriggerUpdate",
			Handler:    _KainotomiaService_TriggerUpdate_Handler,
		},
		{
			MethodName: "DeletePlaylist",
			Handler:    _KainotomiaService_DeletePlaylist_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kainotomia/v1alpha1/kainotomia.proto",
}
