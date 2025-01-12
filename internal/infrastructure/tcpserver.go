package infrastructure

import (
	"fmt"
	"net"
	"os"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

const DEFAULT_MAX_CONNECTIONS = 10
const DEFAULT_MAX_QUEUED_CONNECTIONS = DEFAULT_MAX_CONNECTIONS
const DEFAULT_MAX_QUEUED_JOBS = 36
const DEFAULT_ADDRESS = "localhost"
const DEFAULT_PORT = 10000

type TCPServer struct {
	jobQueue          chan *Job
	connectionsQueue  chan net.Conn
	maxConnections    uint
	activeConnections map[string]*Connection
	address           string
	port              uint
}

type TCPServerConfig struct {
	maxConnections       uint
	maxQueuedJobs        uint
	maxQueuedConnections uint
	address              string
	port                 uint
}

func NewDefaultTCPServerConfig() *TCPServerConfig {
	return &TCPServerConfig{
		maxConnections:       DEFAULT_MAX_CONNECTIONS,
		maxQueuedConnections: DEFAULT_MAX_QUEUED_CONNECTIONS,
		maxQueuedJobs:        DEFAULT_MAX_QUEUED_JOBS,
		address:              DEFAULT_ADDRESS,
		port:                 DEFAULT_PORT,
	}
}

func NewTCPServer(config *TCPServerConfig) *TCPServer {
	if config == nil || config.maxQueuedConnections == 0 || config.maxQueuedJobs == 0 {
		return nil
	}

	return &TCPServer{
		jobQueue:          make(chan *Job, config.maxQueuedJobs),
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

		case job := <-server.jobQueue:
			err := server.handleJob(job)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (server *TCPServer) PrepareDisk(disk *models.Disk) error {
	prepareDiskJob := PrepareDiskJob{
		DiskSize: disk.GetSpaceLeft(),
	}

	raw := prepareDiskJob.Bytes()

	header := JobHeader{
		Version:  VERSION,
		Opcode:   PrepareDisk,
		Encoding: EncodingNone,
		DataSize: uint16(len(raw)),
	}

	job := new(Job)
	job.Header = header
	job.Data = raw

	server.jobQueue <- job // TODO: handle max cap

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

	conf := NewDefaultConnectionConfig(conn, server.jobQueue)
	connection := NewConnection(conf)
	server.activeConnections[conn.LocalAddr().String()] = connection
	go connection.Read()

	return nil
}

func (server *TCPServer) handleJob(job *Job) error {
	switch job.Header.Opcode {
	case PrepareDisk:
		return server.prepareDisk(job)
	}

	return &ErrUnknownJob{Opcode: uint8(job.Header.Opcode)}
}

func (server *TCPServer) prepareDisk(job *Job) error {
	var prepareDiskJob PrepareDiskJob
	err := prepareDiskJob.FromBytes(job.Data)

	if err != nil {
		return err
	}

	err = os.Mkdir(os.Getenv("SDISK_HOME")+"/users", 0777)
	if err != nil {
		panic(err)
	}

	return nil
}
