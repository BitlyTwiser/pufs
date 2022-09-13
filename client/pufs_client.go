package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"log"

	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
  "github.com/BitlyTwiser/tinycrypt"
)

var (
  serverPort = flag.String("p", "9000", "Server port")
  serverAddr = flag.String("addr", "127.0.0.1", "Server Address")
  encrypt = flag.Bool("e", true, "Encryptes uploaded files. True by default.")
  password = flag.String("pass", "Testing123@", "Password used to encrypt data")
)

func uploadFile(c pufs_pb.IpfsFileSystemClient, ctx context.Context) error {
  file, err := c.UploadFile(ctx, &pufs_pb.UploadFileRequest{})

  if err != nil {
    return err
  }
  
  //No IPFS hash here, that will not be known until we upload on server
  metadata := &pufs_pb.File{
    Filename: "",
    FileSize: 0,
    IpfsHash: "",
    UploadedAt: timestamppb.New(time.Now()),
  }
  //Byte stream can be encrypted here.
  encryptedData, err := tinycrypt.EncryptByteStream(*password, []byte(""))

  if err != nil {
    return err
  }

  file.Send(&pufs_pb.UploadFileRequest{FileData:  *encryptedData, FileMetadata: metadata})

  return nil
  
}

func printFiles(files pufs_pb.IpfsFileSystem_ListFilesClient) {
  for {
    file, err := files.Recv()
    
    if err == io.EOF {
      break
    }

    if err != nil {
      log.Fatalf("Error reading file stream")
      break
    }
  
    fmt.Printf("File: %v", file.Files)
  }
}

//Note: These are exmples of using the functions. 
func main() {
  flag.Parse()

  conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *serverAddr, *serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

  if err != nil {
    log.Fatalf("Error connection to server: %v", err)
  }

  defer conn.Close()

  c := pufs_pb.NewIpfsFileSystemClient(conn)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)

  defer cancel()

  files, err := c.ListFiles(ctx, &pufs_pb.FilesRequest{})

  if err != nil {
    log.Fatalf("Error getting files %v", err)
  }

  printFiles(files)
}
