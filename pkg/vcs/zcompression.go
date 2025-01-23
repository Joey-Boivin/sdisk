package vcs

import (
	"compress/zlib"
	"io"
)

type ZLibAlgorithm struct {
}

func (z *ZLibAlgorithm) Compress(writer io.Writer, uncompressed []byte) (int, error) {
	zlibWriter := zlib.NewWriter(writer)
	defer zlibWriter.Close()
	return zlibWriter.Write(uncompressed)
}

func (z *ZLibAlgorithm) Uncompress(writer io.Writer, reader io.Reader) (int, error) {
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(writer, zlibReader)
	return int(written), err
}
