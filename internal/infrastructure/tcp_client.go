package infrastructure

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"path/filepath"
	"time"
)

type TCPClient struct {
	jobQueue   chan *Job
	address    string
	port       uint
	connection *Connection
	syncPath   string
}

type TCPClientConfig struct {
	maxQueuedJobs uint
	address       string
	port          uint
	syncPath      string
}

func NewDefaultTCPClientConfig() *TCPClientConfig {
	syncPath := os.Getenv("SDISK_HOME") + DEFAULT_CLIENT_FOLDER

	defaultClientConfig := TCPClientConfig{
		maxQueuedJobs: DEFAULT_MAX_QUEUED_CLIENT_JOBS,
		address:       DEFAULT_ADDRESS,
		port:          DEFAULT_PORT,
		syncPath:      syncPath,
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
	}

	connectionConfig := NewDefaultConnectionConfig(net, client.jobQueue)
	client.connection = NewConnection(connectionConfig)

	return &client
}

func (client *TCPClient) Run() {
	//TODO: This will run continuisly eventually, but only runs once for now
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
		//fmt.Printf("read %d bytes from file\n", read)
		//fmt.Printf("read %s\n", string(fileContentBuffer))
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}

		updateJob := UpdateDataJob{
			Total:    uint64(info.Size()),
			Offset:   uint64(read),
			PathLen:  uint64(len(info.Name())),
			Path:     info.Name(),
			FileData: fileContentBuffer,
		}

		raw, err := updateJob.Bytes() //THE BHUG IS HERE

		if err != nil {
			panic(err)
		}

		header := JobHeader{
			DataSize: uint16(len(raw)),//THE BUG IS HERE
			Version:  VERSION,
			Encoding: EncodingNone,
			Opcode:   UpdateData,
		}

		copy(header.id[:], []byte(TEST_ID))

		job := Job{
			Header: header,
			Data:   raw,
		}

		toSend := job.Bytes()
		wrote, err := client.connection.Write(toSend)
		sent += wrote
		fmt.Printf("sent a packet of total length %d and data length %d\n", wrote, len(raw))

		updateJob.Offset += uint64(wrote)

		if err != nil {
			panic(err)
		}

		total -= read
	}

	fmt.Printf("file sent complete: %s\n", info.Name())
}
