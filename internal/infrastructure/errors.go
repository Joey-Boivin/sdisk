package infrastructure

import "fmt"

type ErrUnknownPacket struct {
	Opcode uint8
}

func (e *ErrUnknownPacket) Error() string {
	return fmt.Sprintf("this packet with opcode %d is unknown", e.Opcode)
}

type ErrUnsuportedProtocolVersion struct {
	ReceivedVersion byte
}

func (e *ErrUnsuportedProtocolVersion) Error() string {
	return fmt.Sprintf("the version %d of the protocol is not supported by the current server", e.ReceivedVersion)
}

type ErrUnexpectedHeaderLength struct {
	ReceivedHeaderLength int
}

func (e *ErrUnexpectedHeaderLength) Error() string {
	return fmt.Sprintf("unexpected header length of %d", e.ReceivedHeaderLength)
}

type ErrIncompletePacket struct {
}

func (e *ErrIncompletePacket) Error() string {
	return "incomplete packet"
}

type ErrMaximumClientsReached struct {
	maxClients uint
}

func (e *ErrMaximumClientsReached) Error() string {
	return fmt.Sprintf("maximum number of %d clients reached", e.maxClients)
}

type ErrDataQueueFilled struct {
}

func (e *ErrDataQueueFilled) Error() string {
	return "data queue has filled"
}

type ErrUnexpectedFileState struct {
}

func (e *ErrUnexpectedFileState) Error() string {
	return "was not able to seek to correct position"
}

type ErrUserHasNoDisk struct {
}

func (e *ErrUserHasNoDisk) Error() string {
	return "user has no disk"
}
