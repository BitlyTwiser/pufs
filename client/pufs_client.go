package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"
	"github.com/BitlyTwiser/tinychunk"

	//	"github.com/BitlyTwiser/tinycrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	serverPort = flag.String("p", "9000", "Server port")
	serverAddr = flag.String("addr", "127.0.0.1", "Server Address")
	encrypt    = flag.Bool("e", true, "Encryptes uploaded files. True by default.")
	password   = flag.String("pass", "Testing123@", "Password used to encrypt data")
	id         int64
)

type Command struct {
	uploadFs         *flag.FlagSet
	downloadFs       *flag.FlagSet
	downloadCappedFs *flag.FlagSet
	listFs           *flag.FlagSet
	deleteFs         *flag.FlagSet
	streamFs         *flag.FlagSet
	uploadData       uploadData
	downloadData     downloadData
	deleteData       deleteData
	command          string
}

type uploadData struct {
	path    *string
	encrypt *bool
	pass    *string
}

type downloadData struct {
	name *string
	path *string
}

type deleteData struct {
	name     *string
	ipfsHash *string
}

func PufsClient() *Command {
	c := &Command{
		uploadFs:         flag.NewFlagSet("upload", flag.ContinueOnError),
		downloadFs:       flag.NewFlagSet("download", flag.ContinueOnError),
		downloadCappedFs: flag.NewFlagSet("download-capped", flag.ContinueOnError),
		listFs:           flag.NewFlagSet("list", flag.ContinueOnError),
		deleteFs:         flag.NewFlagSet("delete", flag.ContinueOnError),
		streamFs:         flag.NewFlagSet("stream", flag.ContinueOnError),
	}

	// Upload Data
	ud := uploadData{
		path:    c.uploadFs.String("path", "", "Path of file to upload"),
		encrypt: c.uploadFs.Bool("encrypt", false, "To encrypt file on upload"),
		pass:    c.uploadFs.String("pass", "", "Password used to encrypt file"),
	}

	c.uploadData = ud

	//Download Data
	dd := downloadData{
		name: c.downloadFs.String("name", "", "Name of file to download"),
		path: c.downloadFs.String("path", "", "Location to download file too"),
	}

	c.downloadData = dd

	// Delete flag
	d := deleteData{
		name: c.deleteFs.String("name", "", "Name of the file to delete"),
	}

	c.deleteData = d

	return c
}

func uploadFileStream(client pufs_pb.IpfsFileSystemClient, fileData *os.File, fileSize int64, fileName string) error {
	var wg sync.WaitGroup
	log.Printf("Sending large file.. File Size: %v", fileSize)
	// Look to make the time variables depending on file size as well.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileUpload, err := client.UploadFileStream(ctx)

	if err != nil {
		return err
	}

	//No IPFS hash here, that will not be known until we upload on server
	metadata := &pufs_pb.File{
		Filename:   fileName,
		FileSize:   fileSize,
		IpfsHash:   "",
		UploadedAt: timestamppb.New(time.Now()),
	}

	data := make([]byte, fileSize)
	_, err = fileData.Read(data)

	if err != nil {
		return err
	}

	log.Println("Sending first request")
	// Send metadata request first then data.
	m := &pufs_pb.UploadFileStreamRequest{Data: &pufs_pb.UploadFileStreamRequest_FileMetadata{
		FileMetadata: metadata,
	}}

	if err := fileUpload.Send(m); err != nil {
		log.Printf("Error sending first request: %v", err)
	}

	// Total chunks is utilized here to ensure we add enough wait groups.
	// In this particular case, we are chunking the data into 2MB chunks. If alloted amount was altered, we would be forced to revisit this logic.
	totalChunks := uint(math.Floor(float64(fileSize) / float64((2 << 20))))

	wg.Add(int(totalChunks))
	err = tinychunk.Chunk(data, 2, func(chunkedData []byte) error {
		defer wg.Done()

		log.Println("Sending chunked data")
		if err := fileUpload.Send(&pufs_pb.UploadFileStreamRequest{Data: &pufs_pb.UploadFileStreamRequest_FileData{FileData: chunkedData}}); err != nil {
			log.Printf("Error sending file: %v", err)
			return err
		}

		return nil
	})

	wg.Wait()

	if err != nil {
		log.Printf("Error chunking and sending data: %v", err)
		return err
	}

	resp, err := fileUpload.CloseAndRecv()

	if err != nil {
		log.Printf("No response from server")
		return err
	}

	if resp.GetSucessful() {
		log.Println("File has been uploaded")
	} else {
		return errors.New("Server did not say successful")
	}

	return nil
}

func uploadFile(path, fileName string, client pufs_pb.IpfsFileSystemClient, ctx context.Context) error {
	file, err := os.OpenFile(fmt.Sprintf("%v%v", path, fileName), os.O_RDONLY, 0400)

	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	//gRPC data size cap at 4MB
	if fileSize >= (2 << 21) {
		log.Println("Sending big file")
		err = uploadFileStream(client, file, fileSize, fileName)

		if err != nil {
			return err
		}
	} else {
		log.Printf("Sending file of size: %v", fileSize)

		fileData := make([]byte, fileSize)
		_, err := file.Read(fileData)

		if err != nil {
			return err
		}

		err = uploadFileData(ctx, client, fileData, fileSize, fileName)

		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFile(fileName string, client pufs_pb.IpfsFileSystemClient, ctx context.Context) error {
	resp, err := client.DeleteFile(ctx, &pufs_pb.DeleteFileRequest{FileName: fileName})

	if err != nil {
		return err
	}

	if resp.Successful {
		log.Println("File deleted")
	} else {
		return errors.New(fmt.Sprintf("Error occured deleting file: %v", resp))
	}

	return nil
}

func downloadCappedFile(fileName, path string, client pufs_pb.IpfsFileSystemClient) error {
	log.Printf("Downloading larger file: %v", fileName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pufs_pb.DownloadFileRequest{FileName: fileName}

	download, err := client.DownloadFile(ctx, req)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("%v/%v", path, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		log.Printf("error opening file to store downloaded data: %v", err)
	}

	for {
		fileChunk, err := download.Recv()

		if err == io.EOF {
			log.Printf("All data downloaded")

			break
		}

		if err != nil {
			log.Printf("Error downloading capped file: %v", err)
			return err
		}

		n, err := file.Write(fileChunk.GetFileData())

		if err != nil {
			return err
		}

		if n == 0 {
			return errors.New("No bytes were written to file!")
		}
	}

	return nil
}

// We must chunk the file here if its over the 4MB limit.
func downloadFile(fileName, path string, client pufs_pb.IpfsFileSystemClient, ctx context.Context) error {
	log.Printf("Downliading file: %v", fileName)

	fileResp, err := client.DownloadUncappedFile(ctx, &pufs_pb.DownloadFileRequest{FileName: fileName})

	if err != nil {
		return err
	}

	fileData, fileMetadata := fileResp.FileData, fileResp.FileMetadata

	log.Println("Downloading file and saving to disk...")

	err = os.WriteFile(fmt.Sprintf("%v/%v", path, fileMetadata.Filename), fileData, 0600)

	if err != nil {
		return err
	}

	return nil
}

//Uploads a file stream that is under the 4MB gRPC file size cap
func uploadFileData(ctx context.Context, client pufs_pb.IpfsFileSystemClient, fileData []byte, fileSize int64, fileName string) error {
	file := &pufs_pb.File{
		Filename:   fileName,
		FileSize:   fileSize,
		IpfsHash:   "",
		UploadedAt: timestamppb.New(time.Now()),
	}

	log.Println("Uploading file")

	request := &pufs_pb.UploadFileRequest{FileData: fileData, FileMetadata: file}
	resp, err := client.UploadFile(ctx, request)

	if err != nil {
		return err
	}

	if !resp.Sucessful {
		return errors.New("Something went wrong uploading file.")
	}

	return nil
}

func printFiles(files pufs_pb.IpfsFileSystem_ListFilesClient) {
	for {
		file, err := files.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading file stream. Error: %v", err)
			break
		}

		fmt.Printf("File: %v", file.Files)
	}
}

// Listen for file changes realtime.
// Take ID and store this upstream.
func subscribeFileStream(client pufs_pb.IpfsFileSystemClient, ctx context.Context) {
	// Call function to remove the client from the subscription when client dies.
	// Note:  This does NOT run. This would work if the client exited gracefully, however, we do not
	// This is effectively an endless loop that checks for streams, so the only way to exist is to liste for an os exit event.
	defer client.UnsubscribeFileStream(ctx, &pufs_pb.FilesRequest{Id: id})

	for {
		stream, err := client.ListFilesEventStream(ctx, &pufs_pb.FilesRequest{Id: id})

		// Retrun on failure
		if err != nil || stream == nil {
			log.Println("Error or stream not empty")
			time.Sleep(time.Second * 5)

			continue
		} else {
			for {
				file, err := stream.Recv()

				if err == io.EOF {
					log.Println("All files read, awaiting..")
					stream = nil
					break
				}

				if err != nil {
					log.Println("error encrountered, retrying..")
					stream = nil
					break
				}

				fmt.Printf("File: %v", file.Files)
			}
		}
	}
}

func chunkFile(fileName string, client pufs_pb.IpfsFileSystemClient) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	size, err := client.FileSize(ctx, &pufs_pb.FileSizeRequest{FileName: fileName})

	if err != nil {
		log.Printf("Could not get file size. Error: %v", err)
	}

	return size.FileSize >= (2 << 20)
}

//Note: These are exmples of using the functions.
func main() {
	flag.Parse()

	// Set client ID
	rand.Seed(time.Now().UTC().UnixNano())

	// Allow for up to 100 clients.
	// Allow this to be command line/config item.
	id = int64(rand.Intn(100))

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *serverAddr, *serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error connection to server: %v", err)
	}

	defer conn.Close()

	c := pufs_pb.NewIpfsFileSystemClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) < 2 {
		fmt.Println("Expected arguments")
		os.Exit(1)
	}

	command := PufsClient()

	switch os.Args[1] {

	case "upload":
		log.Println("Uploading File")
		command.uploadFs.Parse(os.Args[2:])

		if *command.uploadData.path == "" {
			log.Println("Must give path")
			os.Exit(1)
		}

		path, fileName := path.Split(*command.uploadData.path)

		err = uploadFile(path, fileName, c, ctx)

		if err != nil {
			log.Fatalf("Error occured while uploading file: %v", err)
		}
	case "download":
		log.Println("Downloading File")
		command.downloadFs.Parse(os.Args[2:])

		if chunkFile(*command.downloadData.name, c) {
			err = downloadCappedFile(*command.downloadData.name, *command.downloadData.path, c)
		} else {
			err = downloadFile(*command.downloadData.name, *command.downloadData.path, c, ctx)
		}

		if err != nil {
			log.Fatalf("Error downloading file. Error: %v", err)
		}
	case "download-capped":
		// Note: Not really to be used, the above method should take priority in usage as it will split depending on filesize.
		log.Println("Downloading File")
		command.downloadFs.Parse(os.Args[2:])

		err = downloadCappedFile(*command.downloadData.name, *command.downloadData.path, c)
		if err != nil {
			log.Fatalf("Error downloading large file. Error: %v", err)
		}
	case "delete":
		command.deleteFs.Parse(os.Args[2:])
		log.Printf("Deleting File: %v", *command.deleteData.name)

		// Could also delete by IPFS hash
		err = deleteFile(*command.deleteData.name, c, ctx)

		if err != nil {
			log.Fatalf("An error occured whilst deleting file: %v", err)
		}
	case "stream":
		command.streamFs.Parse(os.Args[2:])
		log.Printf("Streaming files..")

		subscribeFileStream(c, ctx)
	case "list":
		log.Println("Listing files")
		command.listFs.Parse(os.Args[2:])

		files, err := c.ListFiles(ctx, &pufs_pb.FilesRequest{})

		if err != nil {
			log.Fatalf("Error getting files %v", err)
		}

		printFiles(files)

	case "default":
		//Print help here
		log.Println("Well nothing happened")
	}
}
