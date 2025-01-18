package infrastructure

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type TCPClient struct {
	jobQueue   chan *Job
	address    string
	port       uint
	connection *Connection
	syncPath   string
	userID     models.UserID
}

type TCPClientConfig struct {
	maxQueuedJobs uint
	address       string
	port          uint
	syncPath      string
	userID        models.UserID
}

func NewDefaultTCPClientConfig(userID models.UserID, host string, port uint, clientRootFolder string) *TCPClientConfig {
	syncPath := os.Getenv("SDISK_HOME") + "/" + clientRootFolder

	defaultClientConfig := TCPClientConfig{
		maxQueuedJobs: DEFAULT_MAX_QUEUED_CLIENT_JOBS,
		address:       host,
		port:          port,
		syncPath:      syncPath,
		userID:        userID,
	}

	return &defaultClientConfig
}

func NewTCPClient(config *TCPClientConfig) *TCPClient {
	if config == nil || config.maxQueuedJobs == 0 {
		return nil
	}

	net, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.address, config.port))

	if err != nil {
		panic(err)
	}

	client := TCPClient{
		jobQueue: make(chan *Job, config.maxQueuedJobs),
		address:  config.address,
		port:     config.port,
		syncPath: config.syncPath,
		userID:   config.userID,
	}

	connectionConfig := NewDefaultConnectionConfig(net, client.jobQueue)
	client.connection = NewConnection(connectionConfig)

	return &client
}

func (client *TCPClient) Run() {
	files := client.walkDirectory()

	for _, file := range files {
		client.sendFile(&file)
	}

	for {
		time.Sleep(1000)
	}
}

type FileToSend struct {
	path  string
	entry fs.DirEntry
}

func (client *TCPClient) walkDirectory() []FileToSend {
	var files []FileToSend
	_ = filepath.WalkDir(client.syncPath, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			files = append(files, FileToSend{path, d})
		}
		return nil
	})

	return files
}

func (client *TCPClient) sendFile(file *FileToSend) {

	entry := file.entry

	info, err := entry.Info()

	path := strings.TrimPrefix(file.path, client.syncPath)

	if err != nil {
		panic(err)
	}

	f, err := os.Open(file.path)
	if err != nil {
		panic(err)
	}

	chunksSize := DEFAULT_QUEUE_SIZE_BYTES - HEADER_SIZE - 24 - len(info.Name())
	chunks := float64(info.Size()) / float64(chunksSize)
	chunksCeil := int(math.Ceil(chunks))

	sent := 0
	read := 0
	total := int(info.Size())

	for i := 0; i < chunksCeil; i++ {
		fileContentBuffer := make([]byte, int(math.Min(float64(total), float64(chunksSize))))

		if err != nil {
			panic(err)
		}

		read, err = f.Read(fileContentBuffer)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}

		updateJob := UpdateDataJob{
			Total:    uint64(info.Size()),
			Offset:   uint64(read),
			PathLen:  uint64(len(path)),
			Path:     path,
			FileData: fileContentBuffer,
		}

		raw, err := updateJob.Bytes()

		if err != nil {
			panic(err)
		}

		header := JobHeader{
			DataSize: uint16(len(raw)),
			Version:  VERSION,
			Encoding: EncodingNone,
			Opcode:   UpdateData,
		}

		idAsBytes := client.userID.Bytes()
		copy(header.id[:], idAsBytes)

		job := Job{
			Header: header,
			Data:   raw,
		}

		toSend := job.Bytes()
		wrote, err := client.connection.Write(toSend)
		sent += wrote

		updateJob.Offset += uint64(wrote)

		if err != nil {
			panic(err)
		}

		total -= read
	}
}
