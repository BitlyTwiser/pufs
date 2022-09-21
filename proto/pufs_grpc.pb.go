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
	UploadFileStream(ctx context.Context, opts ...grpc.CallOption) (IpfsFileSystem_UploadFileStreamClient, error)
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (IpfsFileSystem_DownloadFileClient, error)
	DownloadUncappedFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error)
	ListFiles(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (IpfsFileSystem_ListFilesClient, error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
	ListFilesEventStream(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (IpfsFileSystem_ListFilesEventStreamClient, error)
	UnsubscribeFileStream(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error)
}

type ipfsFileSystemClient struct {
	cc grpc.ClientConnInterface
}

func NewIpfsFileSystemClient(cc grpc.ClientConnInterface) IpfsFileSystemClient {
	return &ipfsFileSystemClient{cc}
}

func (c *ipfsFileSystemClient) UploadFileStream(ctx context.Context, opts ...grpc.CallOption) (IpfsFileSystem_UploadFileStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &IpfsFileSystem_ServiceDesc.Streams[0], "/pufs.IpfsFileSystem/UploadFileStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &ipfsFileSystemUploadFileStreamClient{stream}
	return x, nil
}

type IpfsFileSystem_UploadFileStreamClient interface {
	Send(*UploadFileRequest) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type ipfsFileSystemUploadFileStreamClient struct {
	grpc.ClientStream
}

func (x *ipfsFileSystemUploadFileStreamClient) Send(m *UploadFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ipfsFileSystemUploadFileStreamClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ipfsFileSystemClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	out := new(UploadFileResponse)
	err := c.cc.Invoke(ctx, "/pufs.IpfsFileSystem/UploadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *ipfsFileSystemClient) DownloadUncappedFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error) {
	out := new(DownloadFileResponse)
	err := c.cc.Invoke(ctx, "/pufs.IpfsFileSystem/DownloadUncappedFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *ipfsFileSystemClient) ListFilesEventStream(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (IpfsFileSystem_ListFilesEventStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &IpfsFileSystem_ServiceDesc.Streams[3], "/pufs.IpfsFileSystem/ListFilesEventStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &ipfsFileSystemListFilesEventStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type IpfsFileSystem_ListFilesEventStreamClient interface {
	Recv() (*FilesResponse, error)
	grpc.ClientStream
}

type ipfsFileSystemListFilesEventStreamClient struct {
	grpc.ClientStream
}

func (x *ipfsFileSystemListFilesEventStreamClient) Recv() (*FilesResponse, error) {
	m := new(FilesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ipfsFileSystemClient) UnsubscribeFileStream(ctx context.Context, in *FilesRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error) {
	out := new(UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/pufs.IpfsFileSystem/UnsubscribeFileStream", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IpfsFileSystemServer is the server API for IpfsFileSystem service.
// All implementations must embed UnimplementedIpfsFileSystemServer
// for forward compatibility
type IpfsFileSystemServer interface {
	UploadFileStream(IpfsFileSystem_UploadFileStreamServer) error
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
	DownloadFile(*DownloadFileRequest, IpfsFileSystem_DownloadFileServer) error
	DownloadUncappedFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error)
	ListFiles(*FilesRequest, IpfsFileSystem_ListFilesServer) error
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	ListFilesEventStream(*FilesRequest, IpfsFileSystem_ListFilesEventStreamServer) error
	UnsubscribeFileStream(context.Context, *FilesRequest) (*UnsubscribeResponse, error)
	mustEmbedUnimplementedIpfsFileSystemServer()
}

// UnimplementedIpfsFileSystemServer must be embedded to have forward compatible implementations.
type UnimplementedIpfsFileSystemServer struct {
}

func (UnimplementedIpfsFileSystemServer) UploadFileStream(IpfsFileSystem_UploadFileStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFileStream not implemented")
}
func (UnimplementedIpfsFileSystemServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) DownloadFile(*DownloadFileRequest, IpfsFileSystem_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) DownloadUncappedFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadUncappedFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) ListFiles(*FilesRequest, IpfsFileSystem_ListFilesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFiles not implemented")
}
func (UnimplementedIpfsFileSystemServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedIpfsFileSystemServer) ListFilesEventStream(*FilesRequest, IpfsFileSystem_ListFilesEventStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFilesEventStream not implemented")
}
func (UnimplementedIpfsFileSystemServer) UnsubscribeFileStream(context.Context, *FilesRequest) (*UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFileStream not implemented")
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

func _IpfsFileSystem_UploadFileStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IpfsFileSystemServer).UploadFileStream(&ipfsFileSystemUploadFileStreamServer{stream})
}

type IpfsFileSystem_UploadFileStreamServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*UploadFileRequest, error)
	grpc.ServerStream
}

type ipfsFileSystemUploadFileStreamServer struct {
	grpc.ServerStream
}

func (x *ipfsFileSystemUploadFileStreamServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ipfsFileSystemUploadFileStreamServer) Recv() (*UploadFileRequest, error) {
	m := new(UploadFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _IpfsFileSystem_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpfsFileSystemServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pufs.IpfsFileSystem/UploadFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpfsFileSystemServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _IpfsFileSystem_DownloadUncappedFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpfsFileSystemServer).DownloadUncappedFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pufs.IpfsFileSystem/DownloadUncappedFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpfsFileSystemServer).DownloadUncappedFile(ctx, req.(*DownloadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _IpfsFileSystem_ListFilesEventStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FilesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IpfsFileSystemServer).ListFilesEventStream(m, &ipfsFileSystemListFilesEventStreamServer{stream})
}

type IpfsFileSystem_ListFilesEventStreamServer interface {
	Send(*FilesResponse) error
	grpc.ServerStream
}

type ipfsFileSystemListFilesEventStreamServer struct {
	grpc.ServerStream
}

func (x *ipfsFileSystemListFilesEventStreamServer) Send(m *FilesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _IpfsFileSystem_UnsubscribeFileStream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpfsFileSystemServer).UnsubscribeFileStream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pufs.IpfsFileSystem/UnsubscribeFileStream",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpfsFileSystemServer).UnsubscribeFileStream(ctx, req.(*FilesRequest))
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
			MethodName: "UploadFile",
			Handler:    _IpfsFileSystem_UploadFile_Handler,
		},
		{
			MethodName: "DownloadUncappedFile",
			Handler:    _IpfsFileSystem_DownloadUncappedFile_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _IpfsFileSystem_DeleteFile_Handler,
		},
		{
			MethodName: "UnsubscribeFileStream",
			Handler:    _IpfsFileSystem_UnsubscribeFileStream_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFileStream",
			Handler:       _IpfsFileSystem_UploadFileStream_Handler,
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
		{
			StreamName:    "ListFilesEventStream",
			Handler:       _IpfsFileSystem_ListFilesEventStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/pufs.proto",
}
