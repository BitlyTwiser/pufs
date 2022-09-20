package ipfs

import (
	"errors"
	"fmt"
)

type FileData struct {
	FileName   string
	FileSize   int64
	IpfsHash   string //There might be a special type we can import here for this
	UploadedAt int64
}

type Node struct {
	Data FileData
	next *Node
	prev *Node
}

type IpfsFiles struct {
	length int
	head   *Node
	tail   *Node
}

// Store directories. potentially in another list?
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
