package mocks

import (
	"io"
)

type HasherMock struct {
	io.Writer
	FnSum         func(b []byte) []byte
	SumCalled     bool
	SumCalledWith []byte
}

func (h *HasherMock) Sum(b []byte) []byte {
	h.SumCalled = true
	h.SumCalledWith = b

	if h.FnSum != nil {
		return h.FnSum(b)
	}

	return []byte{}
}

func (h *HasherMock) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (h *HasherMock) BlockSize() int {
	return 0
}

func (h *HasherMock) Size() int {
	return 0
}

func (h *HasherMock) Reset() {

}

type CompressorMock struct {
	FnCompress                       func(writer io.Writer, uncompressed []byte) (int, error)
	FnCompressCalled                 bool
	FnCompressCalledWithWriter       io.Writer
	FnCompressCalledWithUncompressed []byte

	FnUnCompress                 func(writer io.Writer, reader io.Reader) (int, error)
	FnUnCompressCalled           bool
	FnUnCompressCalledWithWriter io.Writer
	FnUnCompressCalledWithReader io.Reader
}

func (c *CompressorMock) Compress(writer io.Writer, uncompressed []byte) (int, error) {
	c.FnCompressCalled = true
	c.FnCompressCalledWithWriter = writer
	c.FnCompressCalledWithUncompressed = uncompressed

	if c.FnCompress != nil {
		return c.FnCompress(writer, uncompressed)
	}

	return 0, nil
}

func (c *CompressorMock) Uncompress(writer io.Writer, reader io.Reader) (int, error) {
	c.FnUnCompressCalled = true
	c.FnUnCompressCalledWithWriter = writer
	c.FnUnCompressCalledWithReader = reader

	if c.FnUnCompress != nil {
		return c.FnUnCompress(writer, reader)
	}

	return 0, nil
}
