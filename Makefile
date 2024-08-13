
BIN_NAME=cdisk

all: format build

format:
	go fmt github.com/Joey-Boivin/cdisk/...

run: build
	bin/${BIN_NAME}

build: prep
	go build -v -o bin/${BIN_NAME} cmd/main.go

test:
	go test -v github.com/Joey-Boivin/cdisk/...

prep:
	mkdir -p bin

clean: 
	rm -rf bin
