
SERVER_BIN_NAME=sdisk-server
CLIENT_BIN_NAME=sdisk-client

all: prep build

prep:
	mkdir -p bin

format:
	go fmt github.com/Joey-Boivin/sdisk/...

lint:
	golangci-lint run ./...

test:
	go test github.com/Joey-Boivin/sdisk/...

build: prep
	go build -o bin/${SERVER_BIN_NAME} cmd/server/main.go
	go build -o bin/${CLIENT_BIN_NAME} cmd/client/main.go

run: build
	bin/${BIN_NAME}

clean: 
	rm -rf bin
