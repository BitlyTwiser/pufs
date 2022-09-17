package ipfs

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	corepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/pkg/errors"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
)

//tie things to struct, create API 1x on creation.
type IpfsNode struct {
	api         icore.CoreAPI
	ctx         context.Context
	LocalFolder string
}

type FileDownloadData struct {
	IpfsHash string
	FileData *[]byte
}

//Return the context from this function, then defer the cancel request until whichever calling function exits
func (i *IpfsNode) Init() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	api, err := spawnNodeFromDefaultFolder(ctx)

	if err != nil {
		panic(`Something went wrong creating the node! 
      This mostly happens when you already have a node running, 
      check if there is an existing node then run the server again`)
	}

	i.api = api
	i.ctx = ctx

	return cancel
}

func setupPlugins(externalPluginsPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	return core.NewNode(ctx, nodeOptions)
}

var loadPluginsOnce sync.Once

// Create node from ./ipfs folder config.
func spawnNodeFromDefaultFolder(ctx context.Context) (icore.CoreAPI, error) {
	defaultPath, err := config.PathRoot()

	if err != nil {
		return nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err
	}

	node, err := createNode(ctx, defaultPath)

	if err != nil {
		return nil, err
	}

	api, err := coreapi.NewCoreAPI(node)

	if err != nil {
		return nil, err
	}

	return api, nil
}

// This can be gained from the ipfs_file_list. The node data will have the ipfs file hash that we can use to remove the file.
func (i *IpfsNode) DeleteFile(ipfsFileHash string) error {
	err := i.api.Pin().Rm(i.ctx, corepath.New(ipfsFileHash))

	if err != nil {
		return err
	}

	return nil
}

//This can take a raw stream of bytes. As the gRPC server accumulates bytes via the client, we can write the file to IPFS.
func (i *IpfsNode) UploadFileAndPin(fileStream *[]byte) (string, error) {
	p, err := i.api.Unixfs().Add(i.ctx, files.NewBytesFile(*fileStream))

	if err != nil {
		return "", err
	}

	err = i.api.Pin().Add(i.ctx, p)
	if err != nil {
		return "", err
	}

	return p.String(), nil
}

func (i *IpfsNode) ListFilesFromNode() error {
	files, err := i.api.Pin().Ls(i.ctx)

	if err != nil {
		return err
	}

	//Provide abetter method of showing these.
	// We will eventually want to return these to the user if they so desire to see the raw string values in ipfs
	for link := range files {
		fmt.Println(link.Path())
	}

	return nil
}

func (i *IpfsNode) GetFile(fileName string, filesList IpfsFiles) (*FileDownloadData, error) {
	// First get file from the nodes
	ipfsHash := filesList.FindNodeFromName(fileName)

	if ipfsHash == "" {
		return nil, errors.Errorf("No file has been found with name: %v", fileName)
	}

	p, err := i.api.Unixfs().Get(i.ctx, corepath.New(ipfsHash))

	if err != nil {
		return nil, err
	}

	size, err := p.Size()

	if err != nil {
		return nil, err
	}

	f := files.ToFile(p)

	fileData := make([]byte, size)
	if _, err := io.ReadFull(f, fileData); err != nil {
		return nil, err
	}

	f.Close()

	return &FileDownloadData{IpfsHash: ipfsHash, FileData: &fileData}, nil
}
