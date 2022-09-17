package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/BitlyTwiser/pufs-server/ipfs"
	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	port        = flag.String("port", "9000", "Set designated server port. Ensure that this will match with client")
	logPath     = flag.String("lp", "./", "Default logging path. Can be overriden.")
	logFileName = flag.String("lfn", "output.log", "Default logging file name. Can be overriden")
)

var (
	logger = log.New(io.MultiWriter(os.Stdout, loggerFile()), "pufs:", log.Llongfile)
)

type IpfsServer struct {
	ipfsNode   ipfs.IpfsNode
	fileSystem ipfs.IpfsFiles
	mutex      sync.RWMutex
	pufs_pb.UnimplementedIpfsFileSystemServer
}

func (i *IpfsServer) StreamFile(stream pufs_pb.IpfsFileSystem_UploadFileStreamServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Println("Uploading File from client")

	for {
		fileData, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			logger.Printf("Error receiving data from client. Error: %v", err)
			return err
		}

		logger.Printf("Uploading file name to IPFS: %v", fileData.FileMetadata.Filename)

		ipfsHash, err := i.ipfsNode.UploadFileAndPin(&fileData.FileData)

		if err != nil {
			logger.Printf("error uploading file to IPFS. Error: %v", err)
			return err
		}

		i.fileSystem.Append(&ipfs.Node{Data: ipfs.FileData{
			FileName:   fileData.FileMetadata.Filename,
			FileSize:   fileData.FileMetadata.FileSize,
			IpfsHash:   ipfsHash,
			UploadedAt: fileData.FileMetadata.UploadedAt.AsTime().Unix(),
		}})

		logger.Println("File added to virtual file system")
	}
}

//For files under the 4MB gRPC file size cap
// DRY up function and remove dupicated code
func (i *IpfsServer) UploadFile(ctx context.Context, fileData *pufs_pb.UploadFileRequest) (*pufs_pb.UploadFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Printf("Uploading file name to IPFS: %v", fileData.FileMetadata.Filename)

	ipfsHash, err := i.ipfsNode.UploadFileAndPin(&fileData.FileData)

	if err != nil {
		logger.Printf("error uploading file to IPFS. Error: %v", err)
		return nil, err
	}

	i.fileSystem.Append(&ipfs.Node{Data: ipfs.FileData{
		FileName:   fileData.FileMetadata.Filename,
		FileSize:   fileData.FileMetadata.FileSize,
		IpfsHash:   ipfsHash,
		UploadedAt: fileData.FileMetadata.UploadedAt.AsTime().Unix(),
	}})

	logger.Println("File added to virtual file system")

	return &pufs_pb.UploadFileResponse{Sucessful: true}, nil
}

// If over the 4MB cap for grpc, split and chunk into multiple files
func (i *IpfsServer) DownloadFile(in *pufs_pb.DownloadFileRequest, stream pufs_pb.IpfsFileSystem_DownloadFileServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	return nil
}

// A simplified download function without streaming for files under the 4MB byte cap forced via gRPC
func (i *IpfsServer) DownloadUncappedFile(ctx context.Context, in *pufs_pb.DownloadFileRequest) (*pufs_pb.DownloadFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

  logger.Printf("Downloading File: %v to %v", in.FileName, i.ipfsNode.LocalFolder)
  
  fileData, err := i.ipfsNode.GetFile(in.FileName, i.fileSystem)

  if err != nil {
    logger.Printf("Error obtaining file: %v", err)
    return nil, err
  }

  logger.Printf("File %v Downloaded", in.FileName)

	return nil, nil
}

func (i *IpfsServer) ListFiles(in *pufs_pb.FilesRequest, stream pufs_pb.IpfsFileSystem_ListFilesServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	files := i.fileSystem.Files()

	logger.Println("Obtaining files")

	for _, f := range files {
		err := stream.Send(&pufs_pb.FilesResponse{Files: &pufs_pb.File{
			Filename:   f.Data.FileName,
			FileSize:   f.Data.FileSize,
			IpfsHash:   f.Data.IpfsHash,
			UploadedAt: timestamppb.New(time.UnixMicro(f.Data.UploadedAt)),
		}})

		if err != nil {
			logger.Printf("Error sending files to client: %v", err)
			return err
		}
	}

	logger.Println("Finished sending files to client")

	return nil
}

func (i *IpfsServer) DeleteFile(ctx context.Context, in *pufs_pb.DeleteFileRequest) (*pufs_pb.DeleteFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	return &pufs_pb.DeleteFileResponse{}, nil
}

func loggerFile() *os.File {
	if _, err := os.Stat(*logPath); os.IsNotExist(err) {
		if err := os.MkdirAll(*logPath, 0700); err != nil {
			panic(fmt.Sprintf("Could not create directory: %v", err))
		}
	}

	f, err := os.OpenFile(*logPath+"/"+*logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		fmt.Println(err)
		panic("Something happened with the logger")
	}

	return f
}

func main() {
	flag.Parse()

	// Setup Ipfs node
	ipfsNode := ipfs.IpfsNode{LocalFolder: "/tmp/"}
	//Setup Ipfs server and return context.
	cancel := ipfsNode.Init()

	defer cancel()

	// Setup virtual File sytem.
	fileSystem := ipfs.IpfsFiles{}

	var opts []grpc.ServerOption
	logger.Printf("Server starting, address: localhost:%v\nLogger started: Logging to path: %v", *port, *logPath)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", *port))
	if err != nil {
		logger.Printf("Error starting listener: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	pufs_pb.RegisterIpfsFileSystemServer(grpcServer, &IpfsServer{ipfsNode: ipfsNode, fileSystem: fileSystem})

	logger.Fatal(grpcServer.Serve(listener))
}
