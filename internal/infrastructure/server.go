package infrastructure

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type JobType uint8

const (
	prepareDiskType JobType = iota
	realTimeData    JobType = iota
	maxQueuedJobs           = 1000
	OS_READ                 = 04
	OS_WRITE                = 02
	OS_EX                   = 01
	OS_USER_SHIFT           = 6
	OS_GROUP_SHIFT          = 3
	OS_OTH_SHIFT            = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

type Job struct {
	size     uint16
	version  uint8
	encoding uint8
	jobType  uint8
	data     []byte
}

type Server struct {
	queue       chan Job
	connections chan net.Conn
}

func NewServer() *Server {
	return &Server{
		queue:       make(chan Job, maxQueuedJobs),
		connections: make(chan net.Conn),
	}
}

func (s *Server) PrepareDisk(d *models.Disk) error {
	var payload bytes.Buffer
	encoder := gob.NewEncoder(&payload)
	size := d.GetSpaceLeft()
	err := encoder.Encode(prepareDiskJob{size})

	if err != nil {
		return err
	}

	job := Job{size: 0, version: 0, encoding: 0, jobType: 0, data: payload.Bytes()}

	s.queue <- job
	return nil
}

func (s *Server) Run() {
	go s.connectionWorker()
	for {
		select {
		case conn := <-s.connections:
			go s.handleClient(conn)

		case job := <-s.queue:
			err := s.handleJob(&job)
			if err != nil {
				fmt.Print(err)
			}
		}
	}
}

func (s *Server) connectionWorker() {
	listener, err := net.Listen("tcp", "localhost:8989")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		s.connections <- conn
	}
}

func (s *Server) handleClient(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		amount, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d bytes\n", amount)
	}
}

func (s *Server) handleJob(job *Job) error {
	switch job.jobType {
	case uint8(prepareDiskType):
		return s.prepareDisk(job)
	}
	return &ErrUnknownJob{Opcode: job.jobType}
}

type prepareDiskJob struct {
	Size uint64
}

func (s *Server) prepareDisk(job *Job) error {
	fmt.Println("Preparing the disk!")
	r := bytes.NewReader(job.data)
	decoder := gob.NewDecoder(r)
	var p prepareDiskJob
	err := decoder.Decode(&p)

	if err != nil {
		return err
	}

	dir_file_mode := os.ModeDir | (OS_USER_RWX)
	err = os.Mkdir("users", dir_file_mode)
	if err != nil {
		panic(err)
	}

	return nil
}
