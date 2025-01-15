package infrastructure

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/smallnest/ringbuffer"
)

const DEFAULT_QUEUE_SIZE_BYTES = 1024 * 16 // 1MiB TODO: Dynamic size buffer depending on the transaction to optmize performance and ram usage

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

/*
func (connection *Connection) Read() {
	for {
		totalRead := 0
		previousPacketLength := 0
		nextPacketLength := 0
		buff := make([]byte, connection.dataQueueSizeBytes)
		var header JobHeader
		var job Job

		for ok := true; ok; ok = (totalRead != nextPacketLength) {
			read, err := connection.conn.Read(buff[totalRead:])
			fmt.Printf("read %d \n", read)
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				return
			}

			if totalRead == 0 {
				header.fromBytes(buff)
				nextPacketLength = int(header.DataSize) + HEADER_SIZE
			}

			totalRead += read

			if totalRead > nextPacketLength {
				fmt.Printf("More than one!\n")
				err = job.FromBytes(buff)
				if err != nil {
					panic(err)
				}
				previousPacketLength = nextPacketLength
				connection.jobQueue <- &job //queue first job
				header.fromBytes(buff[nextPacketLength:])
				totalRead -= nextPacketLength
				nextPacketLength = int(header.DataSize) + HEADER_SIZE
			}
		}

		_ = job.FromBytes(buff[previousPacketLength:])
		connection.jobQueue <- &job
	}
}
*/

func (connection *Connection) Read() {
	ring := ringbuffer.New(DEFAULT_QUEUE_SIZE_BYTES * 14)

	for {
		buff := make([]byte, connection.dataQueueSizeBytes)
		readDeadline := time.Now().Add(100 * time.Millisecond)
		_ = connection.conn.SetDeadline(readDeadline)
		read, err := connection.conn.Read(buff)

		if err != nil && err != io.EOF && !os.IsTimeout(err) {
			panic(err)
		}

		wrote, err := ring.Write(buff[:read])
		if err != nil {
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
			panic(err)
		}

		err = h.fromBytes(buff[:HEADER_SIZE])
		if err != nil {
			panic(err)
		}

		nextPacketLength := int(h.DataSize) + HEADER_SIZE

		if ring.Length() < nextPacketLength {
			continue
		}

		readFromRing, err := ring.Read(buff[:nextPacketLength])
		fmt.Printf("read %d from ring buff\n", readFromRing)
		if err != nil {
			panic(err)
		}

		var job Job
		err = job.FromBytes(buff[:nextPacketLength])
		if err != nil {
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
