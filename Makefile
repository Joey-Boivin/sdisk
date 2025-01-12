BIN_SERVER_NAME=sdisk-server
BIN_CLIENT_NAME=sdisk-client
BOLD_YELLOW=\033[1;33m
BOLD_GREEN=\033[1;32m
COLOR_OFF=\033[0m

INFO_PREP = ${BOLD_GREEN}INFO: Target prep${COLOR_OFF}
INFO_FORMAT = ${BOLD_GREEN}INFO: Target format${COLOR_OFF}
INFO_LINT = ${BOLD_GREEN}INFO: Target lint${COLOR_OFF}
INFO_TEST = ${BOLD_GREEN}INFO: Target test${COLOR_OFF}
INFO_BENCHMARK_FILES = ${BOLD_GREEN}INFO: Target benchmark_files${COLOR_OFF}
INFO_BUILD = ${BOLD_GREEN}INFO: Target build${COLOR_OFF}
INFO_RUN = ${BOLD_GREEN}INFO: Target run${COLOR_OFF}
INFO_CLEAN = ${BOLD_GREEN}INFO: Target clean${COLOR_OFF}
INFO_TEST_FILES =  ${BOLD_YELLOW}WARN: Creating large files for tests in data folder. Don't forget to make clean!${COLOR_OFF}
WARN_BENCHMARK_FILES = ${BOLD_YELLOW}WARN: Creating large files for benchmarks in data folder. Don't forget to make clean!${COLOR_OFF}

all: build

prep:
	@. scripts/env.sh
	@echo -e "${INFO_PREP}"
	@mkdir -p bin
	@mkdir -p "${SDISK_ROOT}"

format:
	@echo -e "${INFO_FORMAT}"
	@go fmt github.com/Joey-Boivin/sdisk/...

lint:
	@echo -e "${INFO_LINT}"
	@golangci-lint run ./...

test:
	@echo -e "${INFO_TEST}"
	@go test github.com/Joey-Boivin/sdisk/...

test_files:
	@echo -e "${INFO_TEST_FILES}"
	@mkdir -p data
	@touch data/greetings.txt && echo "Greetings" > data/greetings.txt

benchmark_files:
	@echo -e "${INFO_BENCHMARK_FILES}"
	@echo -e "${WARN_BENCHMARK_FILES}"
	@mkdir -p data
	@dd if=/dev/random of=data/file1.txt conv=notrunc bs=1024 count=100
	@dd if=/dev/random of=data/file2.txt conv=notrunc bs=1024 count=1000
	@dd if=/dev/random of=data/file3.txt conv=notrunc bs=1024 count=100000
	@dd if=/dev/random of=data/file4.txt conv=notrunc bs=1024 count=1000000

build: prep
	@echo -e "${INFO_BUILD}"
	@go build -o bin/${BIN_SERVER_NAME} cmd/server/main.go
	@go build -o bin/${BIN_CLIENT_NAME} cmd/client/main.go

run: build
	@echo -e "${INFO_RUN}"
	@bin/${BIN_SERVER_NAME}

clean:
	@echo -e "${INFO_CLEAN}"
	@rm -rf bin
	@rm -rf data