package ipfs_test

import (
	"log"
	"os"
	"testing"
	"time"

	. "github.com/BitlyTwiser/pufs-server/ipfs"
	"github.com/stretchr/testify/assert"
)

func TestIPFS(t *testing.T) {

	file, err := os.ReadFile("../assets/testing_files/test.txt")
	assert.Nil(t, err)

	// A localfolder is necessary when creating the IpfsNode element.
	// This is where all files will be stored on the server whilst being retrevied and uploaded.
	// This is a temprary location, the files are then removed.
	//This can be anywhere, but tmp is a decent solution as that is removed
	inode := IpfsNode{LocalFolder: "/tmp/"}

	//Handle context cancellation that bubbles up from init function.
	cancel := inode.Init()

	defer cancel()

	ipfsFileHash, err := inode.UploadFileAndPin(file)
	assert.Nil(t, err)

	t.Logf("File uploaded and pinned. File Hash: %v", ipfsFileHash)

	//Insert node into file list and display data.
	fileSystem := IpfsFiles{}

	t.Log("Printing IPFS File system data")
	fileSystem.Append(&Node{Data: FileData{"test.txt", 0, ipfsFileHash, time.Now().Unix()}})

	t.Logf("Removing Head File From File System: %v", ipfsFileHash)
	fileSystem.Remove("test.txt")

	fileSystem.Append(&Node{Data: FileData{"anotherTest.txt", 0, ipfsFileHash, time.Now().Unix()}})
	fileSystem.Append(&Node{Data: FileData{"anotherTestFinal.txt", 0, ipfsFileHash, time.Now().Unix()}})

	fileSystem.PrintFileName()
	t.Log(fileSystem.Files())
	t.Logf("Length of current file system: %v", fileSystem.Len())
	//fileSystem.Reverse()

	//The IpfsNode needs an instance of the fileSystem type passed here to access the IpfsList
	t.Log("Getting File")
	f, err := inode.GetFile("anotherTest.txt", fileSystem)

	//Get File
	assert.Nil(t, err)
	assert.NotNil(t, f)

	t.Logf("FileData From getting file: %v. Data: %v", f.IpfsHash, *f.FileData)

	err = inode.ListFilesFromNode()

	assert.Nil(t, err)
	//Should remove all matching files
	t.Logf("Removing File From File System: %v", ipfsFileHash)
	//Takes in a file name, this will be passed via the client
	fileSystem.Remove("anotherTestFinal.txt")
	t.Log("Listing files again after removal")
	fileSystem.PrintFileName()
	t.Logf("Length of File System after removing file: %v", fileSystem.Len())

	t.Logf("Remove Hash: %v", ipfsFileHash)

	t.Logf("Removing file: %v from ipfs node.", ipfsFileHash)
	err = inode.DeleteFile(ipfsFileHash)

	assert.Nil(t, err)

	t.Log("Printing Files hosted on IPFS node after removal")
	err = inode.ListFilesFromNode()
}

// Test dumping filesystem
func TestFileSystemDump(t *testing.T) {
	fileSystem := IpfsFiles{DataPath: "../assets/backup_files/file-system-data.bin"}

	t.Log("Printing IPFS File system data")
	fileSystem.Append(&Node{Data: FileData{"test.txt", 0, "ImmaHash", time.Now().Unix()}})
	fileSystem.Append(&Node{Data: FileData{"anotherTest.txt", 0, "ImmaHash", time.Now().Unix()}})

	err := fileSystem.WriteFileSystemDataToDisk()

	assert.Nil(t, err)
}

// Test filesystem restore
func TestFileSystemRestore(t *testing.T) {
	fileSystem := IpfsFiles{DataPath: "../assets/backup_files/file-system-data.bin"}

	err := fileSystem.LoadFileSystemData()

	assert.Nil(t, err)

	log.Println("File system is not empty after restoration")
	assert.NotEmpty(t, fileSystem.Files())

	log.Println("Printing File system is after restoration")
	for _, v := range fileSystem.Files() {
		log.Println(v.Data.FileName)
	}
}
