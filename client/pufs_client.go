package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
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
)

type Command struct {
	uploadFs     *flag.FlagSet
	downloadFs   *flag.FlagSet
	listFs       *flag.FlagSet
	uploadData   uploadData
	downloadData downloadData
	command      string
}

type uploadData struct {
	name    *string
	path    *string
	encrypt *bool
	pass    *string
}

type downloadData struct {
	name *string
	path *string
}

func PufsClient() *Command {
	c := &Command{
		uploadFs:   flag.NewFlagSet("upload", flag.ContinueOnError),
		downloadFs: flag.NewFlagSet("download", flag.ContinueOnError),
		listFs:     flag.NewFlagSet("list", flag.ContinueOnError),
	}

	// Upload Data
	ud := uploadData{
		name:    c.uploadFs.String("name", "", "Name of file to upload"),
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

	return c
}

func uploadFileStream(ctx context.Context, client pufs_pb.IpfsFileSystemClient, fileData *os.File, fileSize int64) error {
	fileUpload, err := client.UploadFileStream(ctx)

	if err != nil {
		return err
	}

	//No IPFS hash here, that will not be known until we upload on server
	metadata := &pufs_pb.File{
		Filename:   "",
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

func uploadFile(client pufs_pb.IpfsFileSystemClient) error {
	file, err := os.OpenFile("../testing/testing_files/iamatotallydifferentfile.txt", os.O_RDONLY, 0400)

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
		err = uploadFileStream(ctx, client, file, fileSize)

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

		// This shoudl all be tied to custom struct type
		err = uploadFileData(ctx, client, fileData, fileSize)

		if err != nil {
			return err
		}
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

// Chunks large file into 2MB segments.
// 2 MB was selected here as it is a nice even number to segment the files bytes streams into and is below the 4MB limit.
// Pass in function that performs actions on the chunked data?
func fileChunker(fileData []byte, fileSize int64, chunkAction func([]byte) error) error {
	//Calculate the chunks in 2MB segments, make this variable?
	chunkSize := 1 << 21

	totalChunks := uint(math.Floor(float64(fileSize) / float64(chunkSize)))

	for i := uint(0); i < totalChunks; i++ {

		err := chunkAction(fileData[:chunkSize])

		if err != nil {
			return err
		}

		chunkSize = chunkSize * 2
	}

	return nil
}

//Uploads a file stream that is under the 4MB gRPC file size cap
func uploadFileData(ctx context.Context, client pufs_pb.IpfsFileSystemClient, fileData []byte, fileSize int64) error {
	file := &pufs_pb.File{
		Filename:   "afilethatistotallyog.txt",
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

	if len(os.Args) < 2 {
		fmt.Println("Expected arguments")
		os.Exit(1)
	}

	command := PufsClient()

	switch os.Args[1] {

	case "upload":
		log.Println("Uploading FIle")
		command.uploadFs.Parse(os.Args[2:])
		log.Println(command.downloadData.name)

		err = uploadFile(c)

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
