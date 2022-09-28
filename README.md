# pufs
- Pufs is a distributes file system utilizing IPFS, gRPC, and protocol buffers.
- The goal is to utilize the server as a primary mechanism to store uplodaed files in memory, allowing the user to stream file system changes to listening clients. The files are stored in a distributed fashion, utilizing IPFS as the backbone to accomplish this.
- IPFS files are pinned, allowing your stored files to forgo garbage collection:
- IPFS docs:

![Usage](https://pufs-gif-github.s3.us-west-2.amazonaws.com/pufs.gif)

# Usage:
- One can utilize the dockerized version of the application or run the server on a node.
- ```make docker``` will use docker compose to build the application using the subsequent dockerfile.

# Cli Usage:
- You can utilize the client cli to interface with the server.
- Examples:
- Upload file:
```
go run ./client/pufs_client.go upload -path <path_to_file>
```
- Stream file changes:
```
go run ./client/pufs_client.go stream
```
- List files currently within the file system:
```
go run ./client/pufs_client.go list
```
- Download file to local filesystem:
- Provide to the client the path of where you want the file to be stored when it is downloaded to disk and the name of the file.
```
go run ./client/pufs_client.go download -name <file_name> -path /tmp
```

# Changes distributes to clients in real time:
- If you utilize the client in stream mode, you will obtain changes to the file system in real time.
- When any user uploads a new file, the user will see all file changes reflected in their client.
- Example:

# Ipfs:
- This application uses [kubo](https://github.com/ipfs/kubo) for spinning up an IPFS node using a local IPFS installation.
- One must have IPFS installed in order to properly utlize the application.
- The dockerized application installs IPFS within the container.
- You can follow the [install](https://docs.ipfs.io/install/) docs for IPFS to get IPFS setup locally.
