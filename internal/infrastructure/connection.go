package infrastructure

import (
	"fmt"
	"io"
	"log"
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
		readDeadline := time.Now().Add(DEFAULT_READ_TIMEOUT_MS * time.Millisecond)
		_ = connection.conn.SetDeadline(readDeadline)
		read, err := connection.conn.Read(buff)

		if err != nil {
			if err == io.EOF {
				break
			}

			if !os.IsTimeout(err) {
				log.Fatal(err.Error())
			}
		}

		wrote, err := ring.Write(buff[:read])

		if err != nil {
			log.Fatalf("expected to write %d and wrote %d with the following error:\n%s", read, wrote, err)
		}

		if ring.Length() < HEADER_SIZE {
			continue
		}

		var h JobHeader
		peeked, err := ring.Peek(buff[:HEADER_SIZE])
		if err != nil {
			log.Fatalf("expected to peek %d and peeked %d with the following error:\n%s", HEADER_SIZE, peeked, err)
		}

		err = h.fromBytes(buff[:HEADER_SIZE])
		if err != nil {
			log.Fatal(err.Error())
		}

		nextPacketLength := int(h.DataSize) + HEADER_SIZE

		if ring.Length() < nextPacketLength {
			continue
		}

		read, err = ring.Read(buff[:nextPacketLength])

		if err != nil {
			fmt.Printf("expected to read %d and read %d with the following error:\n%s", nextPacketLength, read, err)
			ring.Reset()
			continue
		}

		var job Job
		err = job.FromBytes(buff[:nextPacketLength])
		if err != nil {
			log.Fatal(err.Error())
		}

		connection.jobQueue <- &job
	}
}

func (connection *Connection) Write(data []byte) (int, error) {
	return connection.conn.Write(data)
}
