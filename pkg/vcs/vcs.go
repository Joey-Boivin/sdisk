package vcs

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

const Test = 1

const (
	VcsDirectoryName     = ".vcs"
	ObjectsDirectoryName = "objects"
	RefsDirectoryName    = "refs"
	HeadName             = "HEAD"
	HeadInitialValue     = "ref: refs/head/master\n"
)

type Vcs struct {
	rootDirectoryPath    string
	objectsDirectoryPath string
	refsDirectoryPath    string
	headPath             string
}

func NewVcs(path string) *Vcs {
	root := fmt.Sprintf("%s/%s", path, VcsDirectoryName)
	return &Vcs{
		rootDirectoryPath:    root,
		objectsDirectoryPath: fmt.Sprintf("%s/%s", root, ObjectsDirectoryName),
		refsDirectoryPath:    fmt.Sprintf("%s/%s", root, RefsDirectoryName),
		headPath:             fmt.Sprintf("%s/%s", root, HeadName),
	}
}

func (vcs *Vcs) Init() error {
	_, err := os.Stat(vcs.rootDirectoryPath)

	if err == nil {
		return &ErrVcsAlreadyInitialized{}
	}

	if !os.IsNotExist(err) {
		return err
	}

	return vcs.createEmptyRepository()
}

func (vcs *Vcs) Cat(shasum string, writer io.Writer) error {
	f, err := vcs.findObjectFile(shasum)

	if err != nil {
		return err
	}

	return vcs.printObject(f, writer)
}

func (vcs *Vcs) CreateObject(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	hasher := sha1.New()
	io.Writer.Write(hasher, content)
	sum := hasher.Sum(nil)
	path := fmt.Sprintf("%x", sum)
	dir := path[:2]
	fileName := path[2:]

	err = os.Mkdir(fmt.Sprintf("%s/%s", vcs.objectsDirectoryPath, dir), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s/%s", vcs.objectsDirectoryPath, dir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return err
	}

	defer f.Close()
	zlibWriter := zlib.NewWriter(f)
	toWrite := fmt.Sprintf("blob %d\x00%s", len(content), content)
	n, err := zlibWriter.Write([]byte(toWrite))
	zlibWriter.Close()
	fmt.Printf("%d\n", n)
	return err
}

func (vcs *Vcs) createEmptyRepository() error {
	if err := os.MkdirAll(vcs.objectsDirectoryPath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(vcs.refsDirectoryPath, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(vcs.headPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(HeadInitialValue))

	return err
}

func (vcs *Vcs) findObjectFile(shasum string) (io.Reader, error) {
	if len(shasum) < 3 {
		return nil, os.ErrNotExist
	}

	dir := shasum[0:2]
	fileName := shasum[2:]
	path := fmt.Sprintf("%s/%s/%s", vcs.objectsDirectoryPath, dir, fileName)
	return os.OpenFile(path, os.O_RDONLY, 0644)
}

func (vcs *Vcs) printObject(objectReader io.Reader, writer io.Writer) error {
	zlibReader, err := zlib.NewReader(objectReader)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, zlibReader)
	return err
}

type ErrVcsAlreadyInitialized struct {
}

func (e *ErrVcsAlreadyInitialized) Error() string {
	return "already a vcs directory"
}
