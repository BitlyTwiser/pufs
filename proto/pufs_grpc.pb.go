// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/pufs.proto

package proto

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

// IpfsFileSystemClient is the client API for IpfsFileSystem service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IpfsFileSystemClient interface {
	GetFiles(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (*FilesResponse, error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
}

type ipfsFileSystemClient struct {
	cc grpc.ClientConnInterface
}

func NewIpfsFileSystemClient(cc grpc.ClientConnInterface) IpfsFileSystemClient {
	return &ipfsFileSystemClient{cc}
}

func (c *ipfsFileSystemClient) GetFiles(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (*FilesResponse, error) {
	out := new(FilesResponse)
	err := c.cc.Invoke(ctx, "/pufs.IpfsFileSystem/GetFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipfsFileSystemClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error) {
	out := new(DeleteFileResponse)
	err := c.cc.Invoke(ctx, "/pufs.IpfsFileSystem/DeleteFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IpfsFileSystemServer is the server API for IpfsFileSystem service.
// All implementations must embed UnimplementedIpfsFileSystemServer
// for forward compatibility
type IpfsFileSystemServer interface {
	GetFiles(context.Context, *FilesRequest) (*FilesResponse, error)
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	mustEmbedUnimplementedIpfsFileSystemServer()
}

// UnimplementedIpfsFileSystemServer must be embedded to have forward compatible implementations.
type UnimplementedIpfsFileSystemServer struct {
}

func (UnimplementedIpfsFileSystemServer) GetFiles(context.Context, *FilesRequest) (*FilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedIpfsFileSystemServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) mustEmbedUnimplementedIpfsFileSystemServer() {}

// UnsafeIpfsFileSystemServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IpfsFileSystemServer will
// result in compilation errors.
type UnsafeIpfsFileSystemServer interface {
	mustEmbedUnimplementedIpfsFileSystemServer()
}

func RegisterIpfsFileSystemServer(s grpc.ServiceRegistrar, srv IpfsFileSystemServer) {
	s.RegisterService(&IpfsFileSystem_ServiceDesc, srv)
}

func _IpfsFileSystem_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpfsFileSystemServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pufs.IpfsFileSystem/GetFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpfsFileSystemServer).GetFiles(ctx, req.(*FilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpfsFileSystem_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpfsFileSystemServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pufs.IpfsFileSystem/DeleteFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpfsFileSystemServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IpfsFileSystem_ServiceDesc is the grpc.ServiceDesc for IpfsFileSystem service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IpfsFileSystem_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pufs.IpfsFileSystem",
	HandlerType: (*IpfsFileSystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFiles",
			Handler:    _IpfsFileSystem_GetFiles_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _IpfsFileSystem_DeleteFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/pufs.proto",
}
