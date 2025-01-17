package infrastructure

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/smallnest/ringbuffer"
)

type Connection struct {
	conn               net.Conn
	jobQueue           chan *Job
	dataQueueSizeBytes uint
}

type ConnectionConfig struct {
	conn               net.Conn
	dataQueueSizeBytes uint
	jobQueue           chan *Job
}

func NewDefaultConnectionConfig(conn net.Conn, jobQueue chan *Job) *ConnectionConfig {
	return &ConnectionConfig{
		conn:               conn,
		dataQueueSizeBytes: DEFAULT_QUEUE_SIZE_BYTES,
		jobQueue:           jobQueue,
	}
}

func NewConnection(config *ConnectionConfig) *Connection {
	return &Connection{
		conn:               config.conn,
		jobQueue:           config.jobQueue,
		dataQueueSizeBytes: config.dataQueueSizeBytes,
	}
}

func (connection *Connection) Read() {
	ring := ringbuffer.New(DEFAULT_QUEUE_SIZE_BYTES * 10)

	for {
		buff := make([]byte, connection.dataQueueSizeBytes)
		readDeadline := time.Now().Add(DEFAULT_TIMEOUT_MS * time.Millisecond)
		_ = connection.conn.SetDeadline(readDeadline)
		read, err := connection.conn.Read(buff)

		if err != nil {
			if err == io.EOF {
				// TODO: cleanup disconnect
				break
			}

			if !os.IsTimeout(err) {
				// TODO: error channel? + disconnect
				panic(err)
			}
		}

		wrote, err := ring.Write(buff[:read])
		if err != nil {
			// TODO: error channel? + disconnect
			panic(err)
		}

		fmt.Printf("read %d from socket\n", read)
		fmt.Printf("wrote %d to ring buff\n", wrote)

		if ring.Length() < HEADER_SIZE {
			continue
		}

		var h JobHeader
		_, err = ring.Peek(buff[:HEADER_SIZE])
		if err != nil {
			// TODO: error channel? + disconnect
			panic(err)
		}

		err = h.fromBytes(buff[:HEADER_SIZE])
		if err != nil {
			// TODO: error channel? + disconnect
			panic(err)
		}

		nextPacketLength := int(h.DataSize) + HEADER_SIZE

		if ring.Length() < nextPacketLength {
			// TODO: error channel? + disconnect
			continue
		}

		readFromRing, err := ring.Read(buff[:nextPacketLength])
		fmt.Printf("read %d from ring buff\n", readFromRing)
		if err != nil {
			// TODO: error channel? + disconnect
			panic(err)
		}

		var job Job
		err = job.FromBytes(buff[:nextPacketLength])
		if err != nil {
			// TODO: error channel? + disconnect
			panic(err)
		}
		connection.jobQueue <- &job
	}
}

func (connection *Connection) Write(data []byte) (int, error) {
	return connection.conn.Write(data)
}

func (connection *Connection) Stop() {

}
