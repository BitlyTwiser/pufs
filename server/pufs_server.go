package main

import (
	"flag"
	"io"
  "fmt"
	"log"
	"net"
	"os"

  
	"google.golang.org/grpc"
  pufs_pb "github.com/pufs-server/proto"
  //pufs_pb "github.com/BitlyTwiser/pufs-server/proto";
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
  pufs_pb.UnimplementedIpfsFileSystemServer
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

  log.Printf("Server starting on port: %v\nLogger started: Logging to path: %v", *port, *logPath)

  listener, err := net.Listen("tcp",  *port)
  if err != nil {
    log.Println(err)
  }

  grpcServer := grpc.NewServer()
  pufs_pb.RegisterIpfsFileSystemServer(grpcServer, &IpfsServer{}) 

  log.Fatal(grpcServer.Serve(listener))
}
