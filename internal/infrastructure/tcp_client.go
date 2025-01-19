package infrastructure

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type TCPClient struct {
	transactionQueue chan *Transaction
	address          string
	port             uint
	connection       *Connection
	syncPath         string
	userID           models.UserID
}

type TCPClientConfig struct {
	maxQueuedTransactions uint
	address               string
	port                  uint
	syncPath              string
	userID                models.UserID
}

func NewDefaultTCPClientConfig(userID models.UserID, host string, port uint, clientRootFolder string) *TCPClientConfig {
	syncPath := os.Getenv("SDISK_HOME") + "/" + clientRootFolder

	defaultClientConfig := TCPClientConfig{
		maxQueuedTransactions: DEFAULT_MAX_QUEUED_CLIENT_TRANSACTIONS,
		address:               host,
		port:                  port,
		syncPath:              syncPath,
		userID:                userID,
	}

	return &defaultClientConfig
}

func NewTCPClient(config *TCPClientConfig) *TCPClient {
	if config == nil || config.maxQueuedTransactions == 0 {
		return nil
	}

	net, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.address, config.port))

	if err != nil {
		panic(err)
	}

	client := TCPClient{
		transactionQueue: make(chan *Transaction, config.maxQueuedTransactions),
		address:          config.address,
		port:             config.port,
		syncPath:         config.syncPath,
		userID:           config.userID,
	}

	connectionConfig := NewDefaultConnectionConfig(net, client.transactionQueue)
	client.connection = NewConnection(connectionConfig)

	return &client
}

func (client *TCPClient) Run() {
	go client.connection.Read()
	files := walkDirectory(client.syncPath)

	for _, file := range files {
		sendFile(&file, client.syncPath, client.connection, client.userID)
	}

	header := PacketHeader{
		Version:  VERSION,
		Opcode:   PullData,
		Encoding: EncodingNone,
		DataSize: 0,
	}

	copy(header.id[:], client.userID.Bytes())
	data := []byte{}

	packet := Packet{
		Header:  header,
		Payload: data,
	}
	_, err := client.connection.Write(packet.Bytes())

	if err != nil {
		panic(err)
	}

	for transaction := range client.transactionQueue {
		if transaction.packet.Header.Opcode != UpdateData {
			log.Fatalf("unknown opcode")
		}

		err := client.updateData(transaction)

		if err != nil {
			panic(err)
		}
	}
}

func (client *TCPClient) updateData(transaction *Transaction) error {
	var updateDataPayload UpdateDataPayload
	err := updateDataPayload.FromBytes(transaction.packet.Payload)

	if err != nil {
		return err
	}

	filePath := client.syncPath + updateDataPayload.Path
	dirPath := filepath.Dir(filePath)

	err = os.MkdirAll(dirPath, 0777)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if updateDataPayload.Offset != 0 {
		seeked, err := file.Seek(int64(updateDataPayload.Offset), 0)
		if err != nil {
			return err
		}
		if seeked != int64(updateDataPayload.Offset) {
			return &ErrUnexpectedFileState{}
		}
	}

	_, err = file.Write(updateDataPayload.FileData)

	if err != nil {
		return err
	}

	return nil
}
