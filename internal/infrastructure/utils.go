package infrastructure

import (
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/Joey-Boivin/sdisk/internal/models"
)

type FileToSend struct {
	path  string
	entry fs.DirEntry
}

func walkDirectory(dirPath string) []FileToSend {
	var files []FileToSend
	_ = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			files = append(files, FileToSend{path, d})
		}
		return nil
	})

	return files
}

func sendFile(file *FileToSend, syncPath string, connection *Connection, userID models.UserID) {

	entry := file.entry

	info, err := entry.Info()

	path := strings.TrimPrefix(file.path, syncPath)

	if err != nil {
		panic(err)
	}

	f, err := os.Open(file.path)
	if err != nil {
		panic(err)
	}

	chunksSize := DEFAULT_QUEUE_SIZE_BYTES - HEADER_SIZE - 24 - len(info.Name())
	chunks := float64(info.Size()) / float64(chunksSize)
	chunksCeil := int(math.Ceil(chunks))

	sent := 0
	read := 0
	total := int(info.Size())

	for i := 0; i < chunksCeil; i++ {
		fileContentBuffer := make([]byte, int(math.Min(float64(total), float64(chunksSize))))

		if err != nil {
			panic(err)
		}

		read, err = f.Read(fileContentBuffer)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}

		updateJob := UpdateDataJob{
			Total:    uint64(info.Size()),
			Offset:   uint64(read),
			PathLen:  uint64(len(path)),
			Path:     path,
			FileData: fileContentBuffer,
		}

		raw, err := updateJob.Bytes()

		if err != nil {
			panic(err)
		}

		header := JobHeader{
			DataSize: uint16(len(raw)),
			Version:  VERSION,
			Encoding: EncodingNone,
			Opcode:   UpdateData,
		}

		idAsBytes := userID.Bytes()
		copy(header.id[:], idAsBytes)

		job := Job{
			Header: header,
			Data:   raw,
		}

		toSend := job.Bytes()
		wrote, err := connection.Write(toSend)
		sent += wrote

		updateJob.Offset += uint64(wrote)

		if err != nil {
			panic(err)
		}

		total -= read
	}
}
