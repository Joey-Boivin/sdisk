package vcs_test

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"testing"

	"github.com/Joey-Boivin/sdisk/pkg/vcs"
)

const repositoryLocation = "."

func teardown() {
	_ = os.RemoveAll(vcs.VcsDirectoryName)
}

func TestInitVersionControl(t *testing.T) {
	v := vcs.NewVcs(repositoryLocation)
	headPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.HeadName)
	refsPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.RefsDirectoryName)
	objectsPath := fmt.Sprintf("%s/%s", vcs.VcsDirectoryName, vcs.ObjectsDirectoryName)

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
	vcs := vcs.NewVcs(".")

	t.Run("GivenObjectFileExists_WhenCat_ThenPrintWithoutErrors", func(t *testing.T) {
		defer teardown()
		buff := bytes.Buffer{}
		shasum := "0c06f3d6bb103b054c3e8472e95fe6efd74b88b3"
		_ = createFakeObject(shasum)

		err := vcs.Cat(shasum, &buff)

		assertNoError(t, err)
	})

	t.Run("WhenCat_ThenPrintCorrectContent", func(t *testing.T) {
		defer teardown()
		shasum := "0c06f3d6bb103b054c3e8472e95fe6efd74b88b3"
		buff := bytes.Buffer{}
		content := createFakeObject(shasum)

		_ = vcs.Cat(shasum, &buff)

		assertStringEquals(t, buff.String(), content)
	})

	t.Run("GivenObjectFileDoesNotExist_WhenCat_ThenReturnError", func(t *testing.T) {
		defer teardown()
		buff := bytes.Buffer{}
		err := vcs.Cat("nonexistantfileshasum", &buff)

		assertError(t, err)
	})

}

func createFakeObject(shasum string) string {
	content := "hello, world\n"
	dir := ".vcs/objects/" + shasum[:2]
	_ = os.MkdirAll(dir, os.ModePerm)
	f, _ := os.OpenFile(dir+"/"+shasum[2:], os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	w := zlib.NewWriter(f)
	_, _ = w.Write([]byte(content))
	w.Close()
	return content
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

	if got != want {
		t.Fatalf("initial head content was not set properly")
	}
}

func assertStringEquals(t *testing.T, got string, want string) {
	if got != want {
		t.Fatalf("got string %s, but want %s", got, want)
	}
}
