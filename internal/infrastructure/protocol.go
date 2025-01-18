package infrastructure

import (
	"encoding/binary"
)

const VERSION = 0
const HEADER_SIZE = 21
const ID_SIZE = 16

type PacketOpcode byte
type PacketEncoding byte

const (
	PrepareDisk PacketOpcode = iota
	UpdateData
	PullData
)

const (
	EncodingNone PacketEncoding = iota
	EncodingUnused1
	EncodingUnused2
	EncodingUnused3
)

type PacketHeader struct {
	Version  byte
	Opcode   PacketOpcode
	Encoding PacketEncoding
	id       [ID_SIZE]byte
	DataSize uint16
}

type Packet struct {
	Header  PacketHeader
	Payload []byte
}

type PrepareDiskPayload struct {
	DiskSize uint64
}

type UpdateDataPayload struct {
	Total    uint64
	Offset   uint64
	PathLen  uint64
	Path     string
	FileData []byte
}

func (header *PacketHeader) fromBytes(data []byte) error {
	if data[0] != VERSION {
		return &ErrUnsuportedProtocolVersion{ReceivedVersion: data[0]}
	}

	header.Version = data[0]
	header.Opcode = PacketOpcode(data[1])
	header.Encoding = PacketEncoding(data[2])
	copy(header.id[:], data[3:ID_SIZE+3])
	header.DataSize = binary.BigEndian.Uint16(data[3+ID_SIZE : 5+ID_SIZE])
	return nil
}

func (packet *Packet) FromBytes(data []byte) error {
	if len(data) < HEADER_SIZE {
		return &ErrUnexpectedHeaderLength{ReceivedHeaderLength: len(data)}
	}

	var header PacketHeader
	err := header.fromBytes(data[:HEADER_SIZE])
	if err != nil {
		return err
	}

	if int(header.DataSize) > len(data[HEADER_SIZE:]) {
		return &ErrIncompletePacket{}
	}

	packet.Header = header
	packet.Payload = data[HEADER_SIZE : HEADER_SIZE+packet.Header.DataSize]

	return nil
}

func (packet *Packet) Bytes() []byte {
	buff := make([]byte, 0, HEADER_SIZE+len(packet.Payload))
	buffLength := make([]byte, 2)
	binary.BigEndian.PutUint16(buffLength, packet.Header.DataSize)

	buff = append(buff, packet.Header.Version)
	buff = append(buff, byte(packet.Header.Opcode))
	buff = append(buff, byte(packet.Header.Encoding))
	buff = append(buff, packet.Header.id[:]...)
	buff = append(buff, buffLength...)

	return append(buff, packet.Payload...)
}

func (p *PrepareDiskPayload) Bytes() []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, p.DiskSize)
	return buff
}

func (p *PrepareDiskPayload) FromBytes(data []byte) error {
	if len(data) < 8 {
		return &ErrIncompletePacket{}
	}

	p.DiskSize = binary.BigEndian.Uint64(data)
	return nil
}

func (u *UpdateDataPayload) Bytes() ([]byte, error) {
	buff := make([]byte, 0, 24+int(u.PathLen)+len(u.FileData))
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

func (u *UpdateDataPayload) FromBytes(data []byte) error {
	u.Total = binary.BigEndian.Uint64(data[0:8])
	u.Offset = binary.BigEndian.Uint64(data[8:16])
	u.PathLen = binary.BigEndian.Uint64(data[16:24])
	u.Path = string(data[24 : u.PathLen+24])
	u.FileData = data[24+u.PathLen:]
	return nil
}
