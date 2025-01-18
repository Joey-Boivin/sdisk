package infrastructure

import (
	"fmt"
	"net"
	"os"
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
	files := walkDirectory(client.syncPath)

	for _, file := range files {
		sendFile(&file, client.syncPath, client.connection, client.userID)
	}

	for {
		time.Sleep(1000)
	}
}
