package infrastructure

import "github.com/smallnest/ringbuffer"

const DEFAULT_QUEUE_SIZE_BYTES = 10240 // 10KiB

type TCPFramer struct {
	ringbuff *ringbuffer.RingBuffer
}

func NewTCPFramer(queueSize uint) *TCPFramer {
	return &TCPFramer{
		ringbuff: ringbuffer.New(int(queueSize)),
	}
}

func (framer *TCPFramer) Parse(data []byte) (*Job, error) {
	//Add up to cap(ring_buff) - len(ring_buff)
	//dataLength := len(data)
	//ringBufferCapacity := framer.ringbuff.Capacity()
	//ringBufferUsedSpace := framer.ringbuff.Length()

	return nil, nil
}
