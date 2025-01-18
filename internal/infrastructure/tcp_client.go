package infrastructure

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type TCPClient struct {
	packet   chan *Packet
	address    string
	port       uint
	connection *Connection
	syncPath   string
	userID     models.UserID
}

type TCPClientConfig struct {
	maxQueuedPackets uint
	address       string
	port          uint
	syncPath      string
	userID        models.UserID
}

func NewDefaultTCPClientConfig(userID models.UserID, host string, port uint, clientRootFolder string) *TCPClientConfig {
	syncPath := os.Getenv("SDISK_HOME") + "/" + clientRootFolder

	defaultClientConfig := TCPClientConfig{
		maxQueuedPackets: DEFAULT_MAX_QUEUED_CLIENT_PACKETS,
		address:       host,
		port:          port,
		syncPath:      syncPath,
		userID:        userID,
	}

	return &defaultClientConfig
}

func NewTCPClient(config *TCPClientConfig) *TCPClient {
	if config == nil || config.maxQueuedPackets == 0 {
		return nil
	}

	net, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.address, config.port))

	if err != nil {
		panic(err)
	}

	client := TCPClient{
		packet: make(chan *Packet, config.maxQueuedPackets),
		address:  config.address,
		port:     config.port,
		syncPath: config.syncPath,
		userID:   config.userID,
	}

	connectionConfig := NewDefaultConnectionConfig(net, client.packet)
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
