// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: postpb/v1/post.proto

package postpbv1

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
	V1PostService_CreatePost_FullMethodName = "/postpb.v1.V1PostService/CreatePost"
)

// V1PostServiceClient is the client API for V1PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type V1PostServiceClient interface {
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error)
}

type v1PostServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewV1PostServiceClient(cc grpc.ClientConnInterface) V1PostServiceClient {
	return &v1PostServiceClient{cc}
}

func (c *v1PostServiceClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error) {
	out := new(CreatePostResponse)
	err := c.cc.Invoke(ctx, V1PostService_CreatePost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// V1PostServiceServer is the server API for V1PostService service.
// All implementations must embed UnimplementedV1PostServiceServer
// for forward compatibility
type V1PostServiceServer interface {
	CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
	mustEmbedUnimplementedV1PostServiceServer()
}

// UnimplementedV1PostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedV1PostServiceServer struct {
}

func (UnimplementedV1PostServiceServer) CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedV1PostServiceServer) mustEmbedUnimplementedV1PostServiceServer() {}

// UnsafeV1PostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to V1PostServiceServer will
// result in compilation errors.
type UnsafeV1PostServiceServer interface {
	mustEmbedUnimplementedV1PostServiceServer()
}

func RegisterV1PostServiceServer(s grpc.ServiceRegistrar, srv V1PostServiceServer) {
	s.RegisterService(&V1PostService_ServiceDesc, srv)
}

func _V1PostService_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(V1PostServiceServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: V1PostService_CreatePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(V1PostServiceServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// V1PostService_ServiceDesc is the grpc.ServiceDesc for V1PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var V1PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "postpb.v1.V1PostService",
	HandlerType: (*V1PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePost",
			Handler:    _V1PostService_CreatePost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "postpb/v1/post.proto",
}
