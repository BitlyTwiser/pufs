package ipfs_test

import (
	. "github.com/BitlyTwiser/pufs-server/ipfs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestIPFS(t *testing.T) {
	file, err := os.ReadFile("../testing/testing_files/test.txt")
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
	fileSystem.Append(&Node{Data: FileData{"anotherTest.txt", 0, ipfsFileHash, time.Now().Unix()}})
	fileSystem.Append(&Node{Data: FileData{"anotherTestFinal.txt", 0, ipfsFileHash, time.Now().Unix()}})

	fileSystem.PrintFileName()
	t.Log(fileSystem.Files())
	t.Logf("Length of current file system: %v", fileSystem.Len())
	//fileSystem.Reverse()

	//The IpfsNode needs an instance of the fileSystem type passed here to access the IpfsList
	t.Log("Getting File")
	f, err := inode.GetFile("test.txt", fileSystem)

	//Get File
	assert.Nil(t, err)
	assert.NotNil(t, f)

	t.Logf("FileData From getting file: %v. Data: %v", f.IpfsHash, *f.FileData)

	err = inode.ListFilesFromNode()

	assert.Nil(t, err)
	//Should remove all matching files
	t.Logf("Removing File From File System: %v", ipfsFileHash)
	//Takes in a file name, this will be passed via the client
	fileSystem.Remove("anotherTest.txt")
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
