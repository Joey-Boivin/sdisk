package main

import (
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
)

type FileToSend struct {
	path  string
	entry fs.DirEntry
}

func getFiles(path string) []FileToSend {
	var files []FileToSend
	_ = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			files = append(files, FileToSend{path, d})
		}
		return nil
	})

	return files
}

func sendFile(file *FileToSend, conn net.Conn) {

	entry := file.entry

	info, err := entry.Info()
	if err != nil {
		panic(err)
	}

	f, err := os.Open(file.path)
	if err != nil {
		panic(err)
	}

	fileContentBuffer := make([]byte, info.Size())

	appendJob := infrastructure.UpdateDataJob{
		Total:    uint64(info.Size()),
		Offset:   0,
		PathLen:  uint64(len(info.Name())),
		Path:     info.Name(),
		FileData: fileContentBuffer,
	}

	read, err := f.Read(fileContentBuffer)

	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	raw, err := appendJob.Bytes()

	if err != nil {
		panic(err)
	}

	header := infrastructure.JobHeader{
		DataSize: uint16(len(raw)),
		Version:  0,
		Encoding: infrastructure.EncodingNone,
		Opcode:   infrastructure.UpdateData,
	}

	job := infrastructure.Job{
		Header: header,
		Data:   raw,
	}

	toSend := job.Bytes()
	wrote, err := conn.Write(toSend)

	appendJob.Offset += uint64(wrote)

	if err != nil {
		panic(err)
	}

	fmt.Printf("read %d and sent %d\n", read, wrote)
	fmt.Printf("for file %s\n", info.Name())
	fmt.Printf("with content %s\n", string(appendJob.FileData))
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		panic(err)
	}

	path := os.Getenv("SDISK_HOME") + "/cmd/client/data"
	files := getFiles(path)

	for _, file := range files {
		sendFile(&file, conn)
	}

	for {
		//wait for transmission to finish
		time.Sleep(1000)
	}
}
