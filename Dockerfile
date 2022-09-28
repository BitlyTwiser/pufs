FROM golang:1.19

WORKDIR /usr/src/pufs-server
RUN wget https://dist.ipfs.tech/kubo/v0.15.0/kubo_v0.15.0_linux-amd64.tar.gz && \
tar -xvzf kubo_v0.15.0_linux-amd64.tar.gz && cd kubo && bash install.sh
RUN ipfs init
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/pufs-server ./server

CMD ["pufs-server"]

