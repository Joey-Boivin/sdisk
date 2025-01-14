package infrastructure

import (
	"net"
)

const DEFAULT_QUEUE_SIZE_BYTES = 10240 // 10KiB

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
	/*
		for {
			totalRead := 0
			nextPacketLength := 0
			buff := make([]byte, connection.dataQueueSizeBytes)
			var header JobHeader
			var job Job

			for ok := true; ok; ok = (totalRead == nextPacketLength) {
				read, err := connection.conn.Read(buff[totalRead:])
				if err != nil {
					panic(err)
				}

				if totalRead == 0 {
					header.fromBytes(buff)
					nextPacketLength = int(header.DataSize) + HEADER_SIZE
				}

				totalRead += read

				if totalRead == nextPacketLength {
					// packet is exactly complete
					err = job.FromBytes(buff)
					if err != nil {
						panic(err)
					}
					connection.jobQueue <- &job
				} else if totalRead > nextPacketLength {
					//copy(previous, buff[wantedLength:])
					err = job.FromBytes(buff)
					if err != nil {
						panic(err)
					}
					connection.jobQueue <- &job
					header.fromBytes(buff[nextPacketLength:])
					totalRead -= nextPacketLength
					nextPacketLength = int(header.DataSize) + HEADER_SIZE
				}
			}
	*/
	for {
		totalRead := 0
		previousPacketLength := 0
		nextPacketLength := 0
		buff := make([]byte, connection.dataQueueSizeBytes)
		var header JobHeader
		var job Job

		for ok := true; ok; ok = (totalRead != nextPacketLength) {
			read, err := connection.conn.Read(buff[totalRead:])
			if err != nil {
				panic(err)
			}

			if totalRead == 0 {
				header.fromBytes(buff)
				nextPacketLength = int(header.DataSize) + HEADER_SIZE
			}

			totalRead += read

			if totalRead > nextPacketLength {
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

		/*
			if err != nil {
				panic(err)
			}

			fmt.Printf("read %d bytes\n", read)

			n := 0
			parsed := n
			var job *Job
			for ok := true; ok; ok = (n == 0) {
				n, err = connection.framer.Parse(buff[n:read])
				parsed += n
				if n > 0 && err != nil {
					job.FromBytes(buff[n:parsed])
					connection.jobQueue <- job
				}
			}
		*/
	}
}

func (connection *Connection) Write(data []byte) (int, error) {
	return connection.conn.Write(data)
}

func (connection *Connection) Stop() {

}
