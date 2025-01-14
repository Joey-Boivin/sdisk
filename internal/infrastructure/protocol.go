package infrastructure

import (
	"encoding/binary"
)

const VERSION = 0
const HEADER_SIZE = 21
const ID_SIZE = 16
const TEST_ID = "1234567891234567" //TODO

type PacketOpcode byte
type PacketEncoding byte

const (
	PrepareDisk PacketOpcode = iota
	UpdateData
)

const (
	EncodingNone PacketEncoding = iota
	EncodingUnused1
	EncodingUnused2
	EncodingUnused3
)

type JobHeader struct {
	Version  byte
	Opcode   PacketOpcode
	Encoding PacketEncoding
	id       [ID_SIZE]byte
	DataSize uint16
}

type Job struct {
	Header JobHeader
	Data   []byte
}

type PrepareDiskJob struct {
	DiskSize uint64
}

type UpdateDataJob struct {
	Total    uint64
	Offset   uint64
	PathLen  uint64
	Path     string
	FileData []byte
}

func (header *JobHeader) fromBytes(data []byte) {
	header.Version = data[0]
	header.Opcode = PacketOpcode(data[1])
	header.Encoding = PacketEncoding(data[2])
	copy(header.id[:], data)
	header.DataSize = binary.BigEndian.Uint16(data[3+ID_SIZE : 5+ID_SIZE])
}

func (j *Job) FromBytes(data []byte) error {
	if len(data) < HEADER_SIZE {
		return &ErrUnexpectedHeaderLength{ReceivedHeaderLength: len(data)}
	}

	if data[0] != VERSION {
		return &ErrUnsuportedProtocolVersion{ReceivedVersion: data[0]}
	}

	var header JobHeader
	header.fromBytes(data[:HEADER_SIZE])

	if int(header.DataSize) > len(data[HEADER_SIZE:]) {
		return &ErrIncompletePacket{}
	}

	j.Header = header
	j.Data = data[HEADER_SIZE:]

	return nil
}

func (j *Job) Bytes() []byte {
	buff := make([]byte, 0, HEADER_SIZE+len(j.Data))
	buffLength := make([]byte, 2)
	binary.BigEndian.PutUint16(buffLength, j.Header.DataSize)

	buff = append(buff, j.Header.Version)
	buff = append(buff, byte(j.Header.Opcode))
	buff = append(buff, byte(j.Header.Encoding))
	buff = append(buff, j.Header.id[:]...)
	buff = append(buff, buffLength...)

	return append(buff, j.Data...)
}

func (p *PrepareDiskJob) Bytes() []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, p.DiskSize)
	return buff
}

func (p *PrepareDiskJob) FromBytes(data []byte) error {
	if len(data) < 8 {
		return &ErrIncompletePacket{}
	}

	p.DiskSize = binary.BigEndian.Uint64(data)
	return nil
}

func (u *UpdateDataJob) Bytes() ([]byte, error) {
	buff := make([]byte, 0, 12+u.PathLen+u.Total)
	buffTotal := make([]byte, 8)
	buffOffset := make([]byte, 8)
	buffPathLen := make([]byte, 8)
	binary.BigEndian.PutUint64(buffTotal, u.Total)
	binary.BigEndian.PutUint64(buffOffset, u.Offset)
	binary.BigEndian.PutUint64(buffPathLen, u.PathLen)

	buff = append(buff, buffTotal...)
	buff = append(buff, buffOffset...)
	buff = append(buff, buffPathLen...)
	buff = append(buff, []byte(u.Path)...)
	return append(buff, u.FileData...), nil
}

func (u *UpdateDataJob) FromBytes(data []byte) error {
	u.Total = binary.BigEndian.Uint64(data[0:8])
	u.Offset = binary.BigEndian.Uint64(data[8:16])
	u.PathLen = binary.BigEndian.Uint64(data[16:24])
	u.Path = string(data[24 : u.PathLen+24])
	u.FileData = data[24+u.PathLen : 24+u.PathLen+u.Total]
	return nil
}
