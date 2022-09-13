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
  port = flag.String("port", "9000", "Set designated server port. Ensure that this will match with client")
  logPath = flag.String("lp", "./", "Default logging path. Can be overriden.")
  logFileName = flag.String("lfn", "output.log", "Default logging file name. Can be overriden")
)

var (
  logger = log.New(io.MultiWriter(os.Stdout, loggerFile()), "pufs:", log.Llongfile)
)

type IpfsServer struct {
  ipfsNode ipfs.IpfsNode
  fileSystem ipfs.IpfsFiles
  mutex sync.RWMutex
  pufs_pb.UnimplementedIpfsFileSystemServer
}

//Here we will need to upload to ipfs
// Capture IPFS hash, store in virtual file system.
// Initial data valyes will be bytestream and the file metadata
func (ipfs *IpfsServer) UploadFile(steam pufs_pb.IpfsFileSystem_UploadFileServer) error {
  ipfs.mutex.Lock()
  defer ipfs.mutex.Unlock()

  return nil 
}

func (ipfs *IpfsServer) DownloadFile(in *pufs_pb.DownloadFileRequest, stream pufs_pb.IpfsFileSystem_DownloadFileServer) error {
  ipfs.mutex.Lock()
  defer ipfs.mutex.Unlock()

  return nil 
}

func (ipfs *IpfsServer) ListFiles(in *pufs_pb.FilesRequest, stream pufs_pb.IpfsFileSystem_ListFilesServer) error {
  ipfs.mutex.Lock()
  defer ipfs.mutex.Unlock()

  files := ipfs.fileSystem.Files()

  logger.Println("Obtaining files")
  
  for _, f := range files {
    err := stream.Send(&pufs_pb.FilesResponse{Files: &pufs_pb.File{
      Filename: f.Data.FileName,
      FileSize: f.Data.FileSize,
      IpfsHash: f.Data.IpfsHash,
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

func (ipfs *IpfsServer) DeleteFile(ctx context.Context, in *pufs_pb.DeleteFileRequest) (*pufs_pb.DeleteFileResponse, error) {
  ipfs.mutex.Lock()
  defer ipfs.mutex.Unlock()

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
    //Do something better here
    panic("Something happened with the logger")
  }

  return f 
}

func main() {
  //Include IPFS package here.
  //On the given calls above, perform the necessary actions.
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

  listener, err := net.Listen("tcp",  fmt.Sprintf("localhost:%v", *port))
  if err != nil {
    logger.Printf("Error starting listener: %v", err)
  }

  grpcServer := grpc.NewServer(opts...)
  pufs_pb.RegisterIpfsFileSystemServer(grpcServer, &IpfsServer{ipfsNode: ipfsNode, fileSystem: fileSystem}) 

  logger.Fatal(grpcServer.Serve(listener))
}
