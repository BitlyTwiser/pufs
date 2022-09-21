package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"
	"github.com/BitlyTwiser/tinycrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	serverPort = flag.String("p", "9000", "Server port")
	serverAddr = flag.String("addr", "127.0.0.1", "Server Address")
	encrypt    = flag.Bool("e", true, "Encryptes uploaded files. True by default.")
	password   = flag.String("pass", "Testing123@", "Password used to encrypt data")
  id int64
)

type Command struct {
	uploadFs     *flag.FlagSet
	downloadFs   *flag.FlagSet
	listFs       *flag.FlagSet
	deleteFs     *flag.FlagSet
	streamFs     *flag.FlagSet
	uploadData   uploadData
	downloadData downloadData
	deleteData   deleteData
	command      string
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
		uploadFs:   flag.NewFlagSet("upload", flag.ContinueOnError),
		downloadFs: flag.NewFlagSet("download", flag.ContinueOnError),
		listFs:     flag.NewFlagSet("list", flag.ContinueOnError),
		deleteFs:   flag.NewFlagSet("delete", flag.ContinueOnError),
		streamFs:   flag.NewFlagSet("stream", flag.ContinueOnError),
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

func uploadFileStream(ctx context.Context, client pufs_pb.IpfsFileSystemClient, fileData *os.File, fileSize int64, fileName string) error {
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
	//Byte stream can be encrypted here.
	encryptedData, err := tinycrypt.EncryptByteStream("Password123", []byte("Something"))

	if err != nil {
		return err
	}

	if err := fileUpload.Send(&pufs_pb.UploadFileRequest{FileData: *encryptedData, FileMetadata: metadata}); err != nil {
		log.Printf("Error sending file: %v", err)
		return err
	}

	return nil
}

func uploadFile(path, fileName string, client pufs_pb.IpfsFileSystemClient) error {
	file, err := os.OpenFile(fmt.Sprintf("%v%v", path, fileName), os.O_RDONLY, 0400)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		err = uploadFileStream(ctx, client, file, fileSize, fileName)

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

func deleteFile(fileName string, client pufs_pb.IpfsFileSystemClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func downloadFile(fileName string, client pufs_pb.IpfsFileSystemClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Printf("Downliading file: %v", fileName)

	fileResp, err := client.DownloadUncappedFile(ctx, &pufs_pb.DownloadFileRequest{FileName: fileName})

	if err != nil {
		return err
	}

	fileData, fileMetadata := fileResp.FileData, fileResp.FileMetadata

	fmt.Println(fileData, fileMetadata)
	log.Println("Downloading file and saving to disk...")

	err = os.WriteFile(fmt.Sprintf("/tmp/%v", fileMetadata.Filename), fileData, 0700)

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
// Client cancellation afte this function exits.
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

//Note: These are exmples of using the functions.
func main() {
	flag.Parse()

  // Set client ID
  rand.Seed(time.Now().UTC().UnixNano())
  
  // Allow for up to 100 clients.
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
		log.Println("Uploading FIle")
		command.uploadFs.Parse(os.Args[2:])

		if *command.uploadData.path == "" {
			log.Println("Must give path")
			os.Exit(1)
		}

		path, fileName := path.Split(*command.uploadData.path)

		//if dir == "" assume "./" ?
		err = uploadFile(path, fileName, c)

		if err != nil {
			panic("Death thing while uploading")
		}
	case "download":
		log.Println("Downloading File")
		command.downloadFs.Parse(os.Args[2:])
		log.Println(*command.downloadData.name)

		err = downloadFile(*command.downloadData.name, c)

		if err != nil {
			panic("Death downloading file")
		}
	case "delete":
		command.deleteFs.Parse(os.Args[2:])
		log.Printf("Deleting File: %v", *command.deleteData.name)

		// Could also delete by IPFS hash
		err = deleteFile(*command.deleteData.name, c)

		if err != nil {
			panic("Death while deleting")
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
