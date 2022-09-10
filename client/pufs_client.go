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
)

var (
  serverPort = flag.String("p", "9000", "Server port")
  serverAddr = flag.String("addr", "127.0.0.1", "Server Address" )
)

func printFiles(files pufs_pb.IpfsFileSystem_ListFilesClient) {
  for {
    file, err := files.Recv()

    if err == io.EOF {
      break
    }

    if err != nil {
      log.Fatalf("Error reading file stream")
    }
  
    log.Printf("File: %v", file.Files)
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
