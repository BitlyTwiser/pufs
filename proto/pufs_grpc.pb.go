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
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (IpfsFileSystem_UploadFileClient, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (IpfsFileSystem_DownloadFileClient, error)
	ListFiles(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (IpfsFileSystem_ListFilesClient, error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
}

type ipfsFileSystemClient struct {
	cc grpc.ClientConnInterface
}

func NewIpfsFileSystemClient(cc grpc.ClientConnInterface) IpfsFileSystemClient {
	return &ipfsFileSystemClient{cc}
}

func (c *ipfsFileSystemClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (IpfsFileSystem_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &IpfsFileSystem_ServiceDesc.Streams[0], "/pufs.IpfsFileSystem/UploadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &ipfsFileSystemUploadFileClient{stream}
	return x, nil
}

type IpfsFileSystem_UploadFileClient interface {
	Send(*UploadFileRequest) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type ipfsFileSystemUploadFileClient struct {
	grpc.ClientStream
}

func (x *ipfsFileSystemUploadFileClient) Send(m *UploadFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ipfsFileSystemUploadFileClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ipfsFileSystemClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (IpfsFileSystem_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &IpfsFileSystem_ServiceDesc.Streams[1], "/pufs.IpfsFileSystem/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &ipfsFileSystemDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type IpfsFileSystem_DownloadFileClient interface {
	Recv() (*DownloadFileResponse, error)
	grpc.ClientStream
}

type ipfsFileSystemDownloadFileClient struct {
	grpc.ClientStream
}

func (x *ipfsFileSystemDownloadFileClient) Recv() (*DownloadFileResponse, error) {
	m := new(DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ipfsFileSystemClient) ListFiles(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (IpfsFileSystem_ListFilesClient, error) {
	stream, err := c.cc.NewStream(ctx, &IpfsFileSystem_ServiceDesc.Streams[2], "/pufs.IpfsFileSystem/ListFiles", opts...)
	if err != nil {
		return nil, err
	}
	x := &ipfsFileSystemListFilesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type IpfsFileSystem_ListFilesClient interface {
	Recv() (*FilesResponse, error)
	grpc.ClientStream
}

type ipfsFileSystemListFilesClient struct {
	grpc.ClientStream
}

func (x *ipfsFileSystemListFilesClient) Recv() (*FilesResponse, error) {
	m := new(FilesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
	UploadFile(IpfsFileSystem_UploadFileServer) error
	DownloadFile(*DownloadFileRequest, IpfsFileSystem_DownloadFileServer) error
	ListFiles(*FilesRequest, IpfsFileSystem_ListFilesServer) error
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	mustEmbedUnimplementedIpfsFileSystemServer()
}

// UnimplementedIpfsFileSystemServer must be embedded to have forward compatible implementations.
type UnimplementedIpfsFileSystemServer struct {
}

func (UnimplementedIpfsFileSystemServer) UploadFile(IpfsFileSystem_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) DownloadFile(*DownloadFileRequest, IpfsFileSystem_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) ListFiles(*FilesRequest, IpfsFileSystem_ListFilesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFiles not implemented")
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

func _IpfsFileSystem_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IpfsFileSystemServer).UploadFile(&ipfsFileSystemUploadFileServer{stream})
}

type IpfsFileSystem_UploadFileServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*UploadFileRequest, error)
	grpc.ServerStream
}

type ipfsFileSystemUploadFileServer struct {
	grpc.ServerStream
}

func (x *ipfsFileSystemUploadFileServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ipfsFileSystemUploadFileServer) Recv() (*UploadFileRequest, error) {
	m := new(UploadFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _IpfsFileSystem_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IpfsFileSystemServer).DownloadFile(m, &ipfsFileSystemDownloadFileServer{stream})
}

type IpfsFileSystem_DownloadFileServer interface {
	Send(*DownloadFileResponse) error
	grpc.ServerStream
}

type ipfsFileSystemDownloadFileServer struct {
	grpc.ServerStream
}

func (x *ipfsFileSystemDownloadFileServer) Send(m *DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _IpfsFileSystem_ListFiles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FilesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IpfsFileSystemServer).ListFiles(m, &ipfsFileSystemListFilesServer{stream})
}

type IpfsFileSystem_ListFilesServer interface {
	Send(*FilesResponse) error
	grpc.ServerStream
}

type ipfsFileSystemListFilesServer struct {
	grpc.ServerStream
}

func (x *ipfsFileSystemListFilesServer) Send(m *FilesResponse) error {
	return x.ServerStream.SendMsg(m)
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
			MethodName: "DeleteFile",
			Handler:    _IpfsFileSystem_DeleteFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _IpfsFileSystem_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _IpfsFileSystem_DownloadFile_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListFiles",
			Handler:       _IpfsFileSystem_ListFiles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/pufs.proto",
}
