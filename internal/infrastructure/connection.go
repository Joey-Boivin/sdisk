package infrastructure

import (
	"fmt"
	"net"
)

type Connection struct {
	conn     net.Conn
	jobQueue chan *Job
}

type ConnectionConfig struct {
	conn          net.Conn
	dataQueueSize uint
	jobQueue      chan *Job
}

func NewDefaultConnectionConfig(conn net.Conn, jobQueue chan *Job) *ConnectionConfig {
	return &ConnectionConfig{
		conn:          conn,
		dataQueueSize: DEFAULT_QUEUE_SIZE_BYTES,
		jobQueue:      jobQueue,
	}
}

func NewConnection(config *ConnectionConfig) *Connection {
	return &Connection{
		conn:     config.conn,
		jobQueue: config.jobQueue,
	}
}

func (connection *Connection) Read() {
	for {
		fmt.Print("Waiting for data")
		buff := make([]byte, 10*1024)
		n, err := connection.conn.Read(buff)
		fmt.Printf("read %d bytes\n", n)
		if err != nil {
			panic(err)
		}
	}
}

func (connection *Connection) Write(data []byte) (int, error) {

	return 0, nil
}

func (connection *Connection) Stop() {

}
