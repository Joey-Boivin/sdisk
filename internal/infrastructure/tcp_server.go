package infrastructure

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type TCPServer struct {
	packetQueue       chan *Packet
	connectionsQueue  chan net.Conn
	maxConnections    uint
	activeConnections map[string]*Connection
	address           string
	port              uint
}

type TCPServerConfig struct {
	maxConnections       uint
	maxQueuedPackets     uint
	maxQueuedConnections uint
	address              string
	port                 uint
}

func NewDefaultTCPServerConfig(host string, port uint) *TCPServerConfig {
	return &TCPServerConfig{
		maxConnections:       DEFAULT_MAX_CONNECTIONS,
		maxQueuedConnections: DEFAULT_MAX_QUEUED_CONNECTIONS,
		maxQueuedPackets:     DEFAULT_MAX_QUEUED_SERVER_PACKETS,
		address:              host,
		port:                 port,
	}
}

func NewTCPServer(config *TCPServerConfig) *TCPServer {
	if config == nil || config.maxQueuedConnections == 0 || config.maxQueuedPackets == 0 {
		return nil
	}

	return &TCPServer{
		packetQueue:       make(chan *Packet, config.maxQueuedPackets),
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

		case packet := <-server.packetQueue:
			err := server.handlePacket(packet)
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

	server.packetQueue <- packet

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

	conf := NewDefaultConnectionConfig(conn, server.packetQueue)
	connection := NewConnection(conf)
	server.activeConnections[conn.LocalAddr().String()] = connection
	go connection.Read()

	return nil
}

func (server *TCPServer) handlePacket(packet *Packet) error {
	switch packet.Header.Opcode {
	case PrepareDisk:
		return server.prepareDisk(packet)
	case UpdateData:
		return server.updateData(packet)
	}

	return &ErrUnknownPacket{Opcode: uint8(packet.Header.Opcode)}
}

func (server *TCPServer) prepareDisk(packet *Packet) error {
	var prepareDiskPayload PrepareDiskPayload
	err := prepareDiskPayload.FromBytes(packet.Payload)

	if err != nil {
		return err
	}

	userID, err := models.FromBytes(packet.Header.id[:])
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

func (server *TCPServer) updateData(packet *Packet) error {
	var updateDataPayload UpdateDataPayload
	err := updateDataPayload.FromBytes(packet.Payload)

	if err != nil {
		return err
	}

	userID, err := models.FromBytes(packet.Header.id[:])
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
