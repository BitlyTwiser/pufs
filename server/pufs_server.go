package main

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/BitlyTwiser/tinychunk"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	port              = flag.String("port", "9000", "Set designated server port. Ensure that this will match with client")
	logPath           = flag.String("lp", "./", "Default logging path. Can be overriden.")
	logFileName       = flag.String("lfn", "output.log", "Default logging file name. Can be overriden")
	fileSystemData    = flag.String("df", "./assets/backup_files/file-system-data.bin", "Default location for data to be written.")
	delFileSystemData = flag.Bool("rd", false, "Delete stored data when server reboots")
)

var (
	logger = log.New(io.MultiWriter(os.Stdout, loggerFile()), "pufs:", log.Llongfile)
)

type IpfsServer struct {
	ipfsNode   ipfs.IpfsNode
	fileSystem ipfs.IpfsFiles
	mutex      sync.RWMutex
	pufs_pb.UnimplementedIpfsFileSystemServer
	fileSub fileSubscriber
}

type fileSubscriber struct {
	eventsChannel chan int
	quitChannel   chan int
	fileEventSubs sync.Map
}

type fileStream struct {
	stream pufs_pb.IpfsFileSystem_ListFilesServer
}

// Used to determine, on the client, to chunk the file download or not.
func (i *IpfsServer) FileSize(ctx context.Context, in *pufs_pb.FileSizeRequest) (*pufs_pb.FileSizeResponse, error) {
	logger.Printf("Obtaioning file size for file: %v", in.FileName)

	i.mutex.Lock()
	defer i.mutex.Unlock()

	d := i.fileSystem.FindNodeDataFromName(in.FileName)

	if d == nil {
		return nil, errors.New(fmt.Sprintf("Error finding file: %v", in.FileName))
	}

	return &pufs_pb.FileSizeResponse{FileSize: d.FileSize}, nil
}

// Look at adding in counter to track amount of bytes written. If the length of written bytes > fileSize, grab until the end of bytestream then stop.
func (i *IpfsServer) UploadFileStream(stream pufs_pb.IpfsFileSystem_UploadFileStreamServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Println("Uploading File stream from client")

	resp, err := stream.Recv()

	if err != nil {
		logger.Printf("Error receiving dataset: %v", err)
		return err
	}

	logger.Printf("Initial requestreceived for file: %v", resp.GetFileMetadata().GetFilename())

	buffer := &bytes.Buffer{}

	for {
		err := contextError(stream.Context())

		if err != nil {
			logger.Printf("Error in context: %v", err)

			return err
		}

		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		logger.Printf("Size of file stream: %v", len(resp.GetFileData()))
		if err != nil {
			logger.Printf("Error receiving data from client. Error: %v", err)
			return err
		}

		buffer.Write(resp.GetFileData())
	}

	logger.Printf("Uploading file name to IPFS: %v", resp.GetFileMetadata().Filename)

	// Store full set of data into IPFS
	ipfsHash, err := i.ipfsNode.UploadFileAndPin(buffer.Bytes())

	if err != nil {
		logger.Printf("error uploading file to IPFS. Error: %v", err)
		return err
	}

	i.fileSystem.Append(&ipfs.Node{Data: ipfs.FileData{
		FileName:   resp.GetFileMetadata().Filename,
		FileSize:   resp.GetFileMetadata().FileSize,
		IpfsHash:   ipfsHash,
		UploadedAt: resp.GetFileMetadata().UploadedAt.AsTime().Unix(),
	}})

	logger.Println("File added to virtual file system")

	// Super hack to avoid blocking on file upload when no receivers.
	if len(i.fileSub.eventsChannel) > 0 {
		_ = <-i.fileSub.eventsChannel
	}

	// Push message that a file was uploaded in case we have subscribers
	i.fileSub.eventsChannel <- 1

	// Store data on local disk
	logger.Println("Saving file to local disk..")
	err = i.fileSystem.WriteFileSystemDataToDisk()

	if err != nil {
		return err
	}

	stream.SendAndClose(&pufs_pb.UploadFileResponse{Sucessful: true})

	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return errors.New("Context cancelled")
	case context.DeadlineExceeded:
		return errors.New("Deadline has been exceded")
	default:
		return nil
	}
}

//For files under the 4MB gRPC file size cap
func (i *IpfsServer) UploadFile(ctx context.Context, fileData *pufs_pb.UploadFileRequest) (*pufs_pb.UploadFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Printf("Uploading file name to IPFS: %v", fileData.FileMetadata.Filename)

	ipfsHash, err := i.ipfsNode.UploadFileAndPin(fileData.FileData)

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

	// Store data on local disk
	logger.Println("Saving file to local disk..")
	err = i.fileSystem.WriteFileSystemDataToDisk()

	if err != nil {
		return nil, err
	}

	// Super hack to avoid blocking on file upload when no receivers.
	if len(i.fileSub.eventsChannel) > 0 {
		_ = <-i.fileSub.eventsChannel
	}

	// Push bool into events channel to force refresh of file clients
	i.fileSub.eventsChannel <- 1

	return &pufs_pb.UploadFileResponse{Sucessful: true}, nil
}

func (i *IpfsServer) getFileMetadata(fileName string) (*ipfs.FileData, error) {
	metadata := i.fileSystem.FindNodeDataFromName(fileName)

	if metadata == nil {
		return nil, errors.New("No file metadata found")
	}

	return metadata, nil
}

// If over the 4MB cap for grpc, split and chunk into multiple files
func (i *IpfsServer) DownloadFile(in *pufs_pb.DownloadFileRequest, stream pufs_pb.IpfsFileSystem_DownloadFileServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Printf("Downloading capped file in chunks..")

	fileData, err := i.ipfsNode.GetFile(in.FileName, i.fileSystem)

	if err != nil {
		return err
	}

	tinychunk.Chunk(*fileData.FileData, 2, func(chunkData []byte) error {
		logger.Printf("File chunk size: %v", len(chunkData))
		d := &pufs_pb.DownloadFileResponseStream{
			Data: &pufs_pb.DownloadFileResponseStream_FileData{FileData: chunkData},
		}

		if err := stream.Send(d); err != nil {
			return err
		}

		return nil
	})

	logger.Println("File downloaded from server")

	return nil
}

// A simplified download function without streaming for files under the 4MB byte cap forced via gRPC
func (i *IpfsServer) DownloadUncappedFile(ctx context.Context, in *pufs_pb.DownloadFileRequest) (*pufs_pb.DownloadFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	fileData, err := i.ipfsNode.GetFile(in.FileName, i.fileSystem)

	if err != nil {
		logger.Printf("Error obtaining file: %v", err)
		return nil, err
	}

	metadata, err := i.getFileMetadata(in.FileName)

	if err != nil {
		return nil, err
	}

	logger.Printf("Downloading File: %v to %v", in.FileName, i.ipfsNode.LocalFolder)

	logger.Printf("File %v Downloaded", in.FileName)
	returnMetadata := &pufs_pb.DownloadFileResponse{
		FileData: *fileData.FileData,
		FileMetadata: &pufs_pb.File{
			Filename:   in.FileName,
			FileSize:   metadata.FileSize,
			IpfsHash:   metadata.IpfsHash,
			UploadedAt: timestamppb.New(time.UnixMicro(metadata.UploadedAt)),
		},
	}

	return returnMetadata, nil
}

func (i *IpfsServer) UnsubscribeFileStream(ctx context.Context, in *pufs_pb.FilesRequest) (*pufs_pb.UnsubscribeResponse, error) {
	logger.Printf("Client with id: %v has disconnected", in.Id)
	i.fileSub.fileEventSubs.Delete(in.Id)

	return &pufs_pb.UnsubscribeResponse{Successful: true}, nil
}

func (i *IpfsServer) ListFilesEventStream(in *pufs_pb.FilesRequest, stream pufs_pb.IpfsFileSystem_ListFilesEventStreamServer) error {
	logger.Printf("Client connected with id: %v", in.Id)
	i.fileSub.fileEventSubs.Store(in.Id, fileStream{stream: stream})

	// Print messages for all subscribers
	for {
		select {
		case <-i.fileSub.eventsChannel:
			logger.Println("Streaming files to clients..")

			i.fileSub.fileEventSubs.Range(func(k, v interface{}) bool {
				// For everystream, sendFiles.
				s, ok := i.fileSub.fileEventSubs.Load(k)

				if !ok {
					logger.Printf("Error filding event from key: %v", k)

					return false
				}

				logger.Printf("Found stream: %v", s)

				stream, ok := s.(fileStream)

				if !ok {
					logger.Println("Error converting stream")

					return false
				}

				err := i.sendFiles(stream.stream)

				if err != nil {
					logger.Printf("Error sending files. Error: %v", err)
					return false
				}

				return true
			})
		case <-i.fileSub.quitChannel:
			logger.Println("Disconnecting")
			return nil
		}
	}
}

func (i *IpfsServer) sendFiles(stream pufs_pb.IpfsFileSystem_ListFilesServer) error {
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
			return err
		}
	}

	return nil
}

func (i *IpfsServer) ListFiles(in *pufs_pb.FilesRequest, stream pufs_pb.IpfsFileSystem_ListFilesServer) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	err := i.sendFiles(stream)

	if err != nil {
		logger.Printf("Error sending files to client: %v", err)
		return err
	}

	logger.Println("Finished sending files to client")

	return nil
}

func (i *IpfsServer) DeleteFile(ctx context.Context, in *pufs_pb.DeleteFileRequest) (*pufs_pb.DeleteFileResponse, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	logger.Printf("Deleting File: %v", in.FileName)
	node := i.fileSystem.FindNode(in.FileName)

	if node == nil {
		logger.Printf("No node with file name: %v was found in filesystem", in.FileName)
		return nil, errors.New("No file found")
	} else {
		i.fileSystem.DeleteNode(node)
	}

	logger.Println("File deleted")

	response := &pufs_pb.DeleteFileResponse{Successful: true}

	return response, nil
}

func loggerFile() *os.File {
	if _, err := os.Stat(*logPath); os.IsNotExist(err) {
		if err := os.MkdirAll(*logPath, 0700); err != nil {
			fmt.Printf("Could not create directory: %v", err)
      os.Exit(1)
		}
	}

	f, err := os.OpenFile(*logPath+"/"+*logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		fmt.Printf("Something happened with the logger")
    os.Exit(1)
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
	fileSystem := ipfs.IpfsFiles{DataPath: *fileSystemData}

	// If the user desires to delete the filesystem and start anew.
	if *delFileSystemData {
		if err := os.Remove(*fileSystemData); err != nil {
			logger.Printf("Error deleting fileSystem data file! You may want to check this manually. Error: %v", err)
		}
	} else {
		logger.Printf("Loading existing file data from %v", *fileSystemData)
		err := fileSystem.LoadFileSystemData()

		if err != nil {
			fmt.Printf("Error loading filesystem data from given location: %v. You may want to check this! Error: %v", *fileSystemData, err)
		}
	}

	var opts []grpc.ServerOption
	logger.Printf("Server starting, address: localhost:%v\nLogger started: Logging to path: %v", *port, *logPath)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		logger.Printf("Error starting listener: %v", err)
	}

	// Create server
	// This unbuffered channel is the reason we have to flush the queues elsewhere.
	eventChannel := make(chan int, 1)

	grpcServer := grpc.NewServer(opts...)
	server := &IpfsServer{
		ipfsNode:   ipfsNode,
		fileSystem: fileSystem,
		fileSub:    fileSubscriber{eventsChannel: eventChannel},
	}

	pufs_pb.RegisterIpfsFileSystemServer(grpcServer, server)

	logger.Fatal(grpcServer.Serve(listener))
}
