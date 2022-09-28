protocomp:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./proto/pufs.proto 

docker: 
	docker compose up -d --build
		
.PHONY: protocomp docker
