
BIN_NAME=sdisk

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
	go build -o bin/${BIN_NAME} cmd/main.go

run: build
	bin/${BIN_NAME}

clean: 
	rm -rf bin
