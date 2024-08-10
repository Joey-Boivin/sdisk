
BIN_NAME=cdisk

all: build

run: build
	bin/${BIN_NAME}

build: prep
	go build -o bin/${BIN_NAME} cmd/main.go

prep:
	mkdir -p bin

clean: 
	rm -rf bin
