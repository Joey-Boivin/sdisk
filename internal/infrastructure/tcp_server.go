package infrastructure

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type Transaction struct {
	packet *Packet
	from   string
}

type TCPServer struct {
	transactionQueue  chan *Transaction
	connectionsQueue  chan net.Conn
	maxConnections    uint
	activeConnections map[string]*Connection
	address           string
	port              uint
}

type TCPServerConfig struct {
	maxConnections        uint
	maxQueuedTransactions uint
	maxQueuedConnections  uint
	address               string
	port                  uint
}

func NewDefaultTCPServerConfig(host string, port uint) *TCPServerConfig {
	return &TCPServerConfig{
		maxConnections:        DEFAULT_MAX_CONNECTIONS,
		maxQueuedConnections:  DEFAULT_MAX_QUEUED_CONNECTIONS,
		maxQueuedTransactions: DEFAULT_MAX_QUEUED_SERVER_TRANSACTIONS,
		address:               host,
		port:                  port,
	}
}

func NewTCPServer(config *TCPServerConfig) *TCPServer {
	if config == nil || config.maxQueuedConnections == 0 || config.maxQueuedTransactions == 0 {
		return nil
	}

	return &TCPServer{
		transactionQueue:  make(chan *Transaction, config.maxQueuedTransactions),
		connectionsQueue:  make(chan net.Conn, config.maxQueuedConnections),
		activeConnections: make(map[string]*Connection),
		maxConnections:    config.maxConnections,
		address:           config.address,
		port:              config.port,
	}
}

func (server *TCPServer) Run() {
	go server.connectionWorker()

	for {
		select {
		case conn := <-server.connectionsQueue:
			err := server.addConnection(conn)
			if err != nil {
				fmt.Println(err)
			}

		case transaction := <-server.transactionQueue:
			err := server.handlePacket(transaction)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (server *TCPServer) PrepareDisk(disk *models.Disk, user *models.User) error {
	prepareDiskPayload := PrepareDiskPayload{
		DiskSize: disk.GetSpaceLeft(),
	}

	raw := prepareDiskPayload.Bytes()

	header := PacketHeader{
		Version:  VERSION,
		Opcode:   PrepareDisk,
		Encoding: EncodingNone,
		DataSize: uint16(len(raw)),
	}

	userID := user.GetID()
	idAsBytes := userID.Bytes()
	copy(header.id[:], idAsBytes)

	packet := new(Packet)
	packet.Header = header
	packet.Payload = raw

	transaction := Transaction{
		packet: packet,
		from:   "localhost",
		//withId: userID.ToString(), TODO?
	}

	server.transactionQueue <- &transaction

	return nil
}

func (server *TCPServer) connectionWorker() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.address, server.port))

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		server.connectionsQueue <- conn
	}
}

func (server *TCPServer) addConnection(conn net.Conn) error {
	if len(server.activeConnections) >= int(server.maxConnections) {
		return &ErrMaximumClientsReached{}
	}

	conf := NewDefaultConnectionConfig(conn, server.transactionQueue)
	connection := NewConnection(conf)
	server.activeConnections[conn.LocalAddr().String()] = connection
	go connection.Read()

	return nil
}

func (server *TCPServer) handlePacket(transaction *Transaction) error {
	switch transaction.packet.Header.Opcode {
	case PrepareDisk:
		return server.prepareDisk(transaction)
	case UpdateData:
		return server.updateData(transaction)
	case PullData:
		return server.pullData(transaction)
	}

	return &ErrUnknownPacket{Opcode: uint8(transaction.packet.Header.Opcode)}
}

func (server *TCPServer) prepareDisk(transaction *Transaction) error {
	var prepareDiskPayload PrepareDiskPayload
	err := prepareDiskPayload.FromBytes(transaction.packet.Payload)

	if err != nil {
		return err
	}

	userID, err := models.FromBytes(transaction.packet.Header.id[:])
	if err != nil {
		return err
	}

	userDiskPath := os.Getenv("SDISK_ROOT") + "/" + userID.ToString()

	err = os.Mkdir(userDiskPath, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (server *TCPServer) updateData(transaction *Transaction) error {
	var updateDataPayload UpdateDataPayload
	err := updateDataPayload.FromBytes(transaction.packet.Payload)

	if err != nil {
		return err
	}

	userID, err := models.FromBytes(transaction.packet.Header.id[:])
	if err != nil {
		return err
	}

	userDiskPath := os.Getenv("SDISK_ROOT") + "/" + userID.ToString()
	info, err := os.Stat(userDiskPath)
	if err != nil || !info.IsDir() {
		return &ErrUserHasNoDisk{}
	}

	filePath := userDiskPath + updateDataPayload.Path
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

func (server *TCPServer) pullData(transaction *Transaction) error {
	userID, err := models.FromBytes(transaction.packet.Header.id[:])
	if err != nil {
		return err
	}

	userDiskPath := os.Getenv("SDISK_ROOT") + "/" + userID.ToString()
	info, err := os.Stat(userDiskPath)
	if err != nil || !info.IsDir() {
		return &ErrUserHasNoDisk{}
	}

	files := walkDirectory(userDiskPath)

	conn := server.activeConnections[transaction.from]
	if conn == nil {
		return &ErrDisconnected{}
	}

	for _, file := range files {
		sendFile(&file, userDiskPath, conn, userID)
	}

	return nil
}
