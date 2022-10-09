package ipfs

import (
	"encoding/binary"
	"encoding/json"
	"io"

	//  "time"
	"errors"
	"fmt"
	"log"
	"os"
)

type FileData struct {
	FileName   string
	FileSize   int64
	IpfsHash   string
	UploadedAt int64
}

// Fix these
type FileDataSerialized struct {
	FileName   string `json:"FileName"`
	FileSize   int64  `json:"FileSize"`
	IpfsHash   string `json:"IpfsHash"`
	UploadedAt int64  `json:"UploadedAt"`
}

type Node struct {
	Data FileData
	next *Node
	prev *Node
}

type IpfsFiles struct {
	length   int
	DataPath string
	head     *Node
	tail     *Node
}

// Store directories. Not used currently.
type IpfsDirectory struct {
	name     string
	contents []IpfsFiles
}

func (f *IpfsFiles) Len() int {
	return f.length
}

func (f *IpfsFiles) Append(node *Node) {
	if f.head == nil {
		f.head = node
		f.tail = node
	} else {
		node.next = f.head
		f.head.prev = node
		f.head = node
	}
	f.length++
}

//Probably will not be used.. could use to get data from a specific file?
func (f IpfsFiles) PrintFileName() {
	for f.head != nil {
		fmt.Printf("%v\n", f.head.Data.FileName)
		f.head = f.head.next
	}
}

//Return array of ipfsfiles within the virtual file system
func (f IpfsFiles) Files() []*Node {
	var files []*Node

	for f.head != nil {
		files = append(files, f.head)
		f.head = f.head.next
	}

	return files
}

func (f IpfsFiles) FirstFile() (FileData, error) {
	if f.head != nil {
		return f.head.Data, nil
	}

	return FileData{}, errors.New("Head value not found")
}

func (f IpfsFiles) LastFile() (FileData, error) {
	if f.tail != nil {
		return f.tail.Data, nil
	}

	return FileData{}, errors.New("Head value not found")
}

func (f IpfsFiles) Reverse() {
	var prev *Node

	currNode := f.head

	f.tail = f.head

	for currNode != nil {
		next := currNode.next
		currNode.next = prev
		prev = currNode
		currNode = next
	}

	f.head = prev

	displayNodeData(f.head)
}

func displayNodeData(node *Node) {
	for node != nil {
		fmt.Printf("%v\n", node.Data.FileName)
		node = node.prev
	}
}

func (f *IpfsFiles) DeleteNode(node *Node) {
	//Head
	if node.prev == nil {
		f.head = node.next
		if f.head == nil {
			f.length--

			return
		}

		node.next.prev = node.prev
		node.prev = node.next

		f.length--

		return
	}

	//Tail
	if node.next == nil {
		node.prev.next = nil
		f.tail = node.prev

		f.length--

		return
	}

	//Middle node
	node.prev.next = node.next
	node.next.prev = node.prev

	node.next = nil
	node.prev = nil

	f.length--
}

func (f IpfsFiles) FindNodeFromName(fileName string) string {
	for f.head != nil {
		if f.head.Data.FileName == fileName {
			return f.head.Data.IpfsHash
		}

		f.head = f.head.next
	}

	return ""
}

func (f IpfsFiles) FindNode(nodeName string) *Node {
	for f.head != nil {
		if f.head.Data.FileName == nodeName {
			return f.head
		}

		f.head = f.head.next
	}

	return nil
}

func (f IpfsFiles) FindNodeDataFromName(fileName string) *FileData {
	for f.head != nil {
		if f.head.Data.FileName == fileName {
			return &f.head.Data
		}

		f.head = f.head.next
	}

	return nil
}

func (f *IpfsFiles) Remove(fileName string) {
	tmp := f.head

	for tmp != nil {
		if tmp.Data.FileName == fileName {
			f.DeleteNode(tmp)
		}

		tmp = tmp.next
	}
}

// Writes file system data to disk
func (f IpfsFiles) WriteFileSystemDataToDisk() error {
	log.Println("Writing filesystem data to disk.")

	file, err := os.OpenFile(f.DataPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	defer file.Close()
	// Write empty file.
	if f.head == nil {
		_, err = file.WriteAt([]byte{}, 0)

		if err != nil {
			return err
		}

		return nil
	}

	for f.head != nil {
		fileBytes, err := json.Marshal(f.head.Data)

		if err != nil {
			return err
		}

		buff := make([]byte, 4)
		binary.LittleEndian.PutUint32(buff, uint32(len(fileBytes)))

		// Write empty buffer between messages.
		_, err = file.Write(buff)
		if err != nil {
			return err
		}

		n, err := file.Write(fileBytes)

		if err != nil {
			return err
		} else if n == 0 {
			return errors.New("zero bytes written to file!")
		}

		f.head = f.head.next
	}

	return nil
}

// Reads file system data from disk
func (f *IpfsFiles) LoadFileSystemData() error {
	var buffData FileDataSerialized

	file, err := os.OpenFile(f.DataPath, os.O_RDONLY, 0600)
	defer file.Close()

	if err != nil {
		return err
	}
	for {
		// Read the delimiter value buffer
		buf := make([]byte, 4)
		_, err := io.ReadFull(file, buf)

		if err == io.EOF {
			break
		}

		size := binary.LittleEndian.Uint32(buf)

		msg := make([]byte, size)
		if _, err := io.ReadFull(file, msg); err != nil {
			return err
		}

		err = json.Unmarshal(msg, &buffData)

		if err != nil {
			return err
		}

		n := &Node{Data: FileData{
			FileName:   buffData.FileName,
			FileSize:   buffData.FileSize,
			IpfsHash:   buffData.IpfsHash,
			UploadedAt: buffData.UploadedAt,
		}}

		f.Append(n)
	}

	return nil
}
