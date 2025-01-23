package vcs_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/Joey-Boivin/sdisk/pkg/vcs"
	"github.com/Joey-Boivin/sdisk/pkg/vcs/mocks"
)

const repositoryLocation = "."

func teardown() {
	_ = os.RemoveAll(vcs.VcsDirectoryName)
}

func TestInitVersionControl(t *testing.T) {
	headPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.HeadName)
	refsPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.RefsDirectoryName)
	objectsPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.ObjectsDirectoryName)

	hasherDummy := mocks.HasherMock{FnSum: func(b []byte) []byte {
		return []byte{}
	}}

	compressorDummy := mocks.CompressorMock{FnCompress: func(writer io.Writer, uncompressed []byte) (int, error) {
		return 0, nil
	}}

	v := vcs.NewVcs(repositoryLocation, &hasherDummy, &compressorDummy)

	t.Run("GivenDirectoryUninitialized_WhenInitializing_CreateRootDirectory", func(t *testing.T) {
		defer teardown()

		err := v.Init()

		assertNoError(t, err)
		assertDirectoryExists(t, vcs.VcsDirectoryName)
	})

	t.Run("GivenDirectoryUninitialized_WhenInitializing_ThenCreateRefsDirectory", func(t *testing.T) {
		defer teardown()

		err := v.Init()

		assertNoError(t, err)
		assertDirectoryExists(t, refsPath)
	})

	t.Run("GivenDirectoryUninitialized_WhenInitializing_ThenCreateObjectsDirectory", func(t *testing.T) {
		defer teardown()

		err := v.Init()

		assertNoError(t, err)
		assertDirectoryExists(t, objectsPath)
	})

	t.Run("GivenDirectoryUninitialized_WhenInitializing_ThenCreateHeadFileCreated", func(t *testing.T) {
		defer teardown()
		err := v.Init()

		assertNoError(t, err)
		assertDirectoryExists(t, headPath)
	})

	t.Run("GivenDirectoryUninitialized_WhenInitializing_ThenHeadHasInitialValue", func(t *testing.T) {
		defer teardown()

		err := v.Init()

		assertNoError(t, err)
		f, _ := os.OpenFile(headPath, os.O_RDONLY, 0644)
		defer f.Close()
		buff := make([]byte, len(vcs.HeadInitialValue))
		_, _ = f.Read(buff)
		assertHeadValue(t, string(buff), vcs.HeadInitialValue)
	})

	t.Run("GivenDirectoryAlreadyInitialized_WhenInitializing_ReturnError", func(t *testing.T) {
		defer teardown()

		_ = v.Init()

		assertError(t, v.Init())
	})
}

func TestCat(t *testing.T) {
	compressedData := "hel"
	unCompressedData := "hello world"

	hasherMock := mocks.HasherMock{FnSum: func(b []byte) []byte {
		return []byte(unCompressedData)
	}}

	compressorMock := mocks.CompressorMock{
		FnCompress: func(writer io.Writer, uncompressed []byte) (int, error) {
			return writer.Write([]byte(compressedData))
		},
		FnUnCompress: func(writer io.Writer, reader io.Reader) (int, error) {
			return writer.Write([]byte(unCompressedData))
		}}

	v := vcs.NewVcs(repositoryLocation, &hasherMock, &compressorMock)

	t.Run("GivenObjectFileExists_WhenCat_ThenPrintWithoutErrors", func(t *testing.T) {
		defer teardown()
		buff := bytes.Buffer{}
		shasum := "0c06f3d6bb103b054c3e8472e95fe6efd74b88b3"
		createFakeCompressedObject(shasum, unCompressedData)

		err := v.Cat(shasum, &buff)

		assertNoError(t, err)
	})

	t.Run("WhenCat_ThenPrintCorrectContent", func(t *testing.T) {
		defer teardown()
		shasum := "0c06f3d6bb103b054c3e8472e95fe6efd74b88b3"
		buff := bytes.Buffer{}
		createFakeCompressedObject(shasum, compressedData)

		_ = v.Cat(shasum, &buff)

		assertStringEquals(t, unCompressedData, buff.String())
	})

	t.Run("GivenObjectFileDoesNotExist_WhenCat_ThenReturnError", func(t *testing.T) {
		defer teardown()
		buff := bytes.Buffer{}
		err := v.Cat("nonexistantfileshasum", &buff)

		assertError(t, err)
	})

}

func TestCreateObject(t *testing.T) {
	compressedData := "hel"
	unCompressedData := "hello world"
	dummyFileName := "text.txt"

	hasherMock := mocks.HasherMock{FnSum: func(b []byte) []byte {
		return []byte(unCompressedData)
	}}

	compressionMock := mocks.CompressorMock{
		FnCompress: func(writer io.Writer, uncompressed []byte) (int, error) {
			return writer.Write([]byte(compressedData))
		},
		FnUnCompress: func(writer io.Writer, reader io.Reader) (int, error) {
			return writer.Write([]byte(unCompressedData))
		}}

	v := vcs.NewVcs(repositoryLocation, &hasherMock, &compressionMock)

	t.Run("GivenFileDoesNotExist_WhenCreateObject_ThenReturnError", func(t *testing.T) {
		defer teardown()
		_ = v.Init()

		_, err := v.CreateObject(dummyFileName)

		assertError(t, err)
	})

	t.Run("GivenFileExists_WhenCreateObject_ThenReturnNoErrors", func(t *testing.T) {
		defer teardown()
		_ = v.Init()
		f := createDummyFile(dummyFileName, []byte(unCompressedData))
		f.Close()
		defer os.Remove(dummyFileName)

		_, err := v.CreateObject(dummyFileName)

		assertNoError(t, err)
	})

	t.Run("WhenCreateObject_ThenCompressedObjectIsInExpectedFormat", func(t *testing.T) {
		defer teardown()
		_ = v.Init()
		f := createDummyFile(dummyFileName, []byte(unCompressedData))
		f.Close()
		defer os.Remove(dummyFileName)

		shasum, _ := v.CreateObject(dummyFileName)

		writer := bytes.Buffer{}
		_ = v.Cat(shasum, &writer)
		assertContentEquals(t, writer.Bytes(), []byte(unCompressedData))
	})
}

func createFakeCompressedObject(shasum string, content string) {
	dir := ".vcs/objects/" + shasum[:2]
	_ = os.MkdirAll(dir, os.ModePerm)
	f, _ := os.OpenFile(dir+"/"+shasum[2:], os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	_, _ = f.Write([]byte(content))
	f.Close()
}

func createDummyFile(path string, content []byte) *os.File {
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	_, _ = f.Write(content)
	_, _ = f.Write(content)
	return f
}

func assertError(t *testing.T, got error) {
	t.Helper()

	if got == nil {
		t.Fatalf("expected an error, but there was no errors")
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf("expected no errors, but got the following error: %s", got.Error())
	}
}

func assertDirectoryExists(t *testing.T, root string) {
	t.Helper()

	_, err := os.Stat(root)
	assertNoError(t, err)
}

func assertHeadValue(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("initial head content was not set properly")
	}
}

func assertStringEquals(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("got string %s, but want %s", got, want)
	}
}

func assertContentEquals(t *testing.T, got []byte, want []byte) {
	t.Helper()

	if !bytes.Contains(got, want) {
		t.Fatalf("content is not equal")
	}
}
