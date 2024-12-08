package models_test

import (
	"testing"

	"github.com/Joey-Boivin/sdisk-api/api/models"
)

func TestAddDisk(t *testing.T) {
	anyUserEmail := "EMAIL@TEST.com"
	anyUserPassword := "12345"
	anyDiskSize := 1024
	user := models.NewUser(anyUserEmail, anyUserPassword)
	anyDisk := models.NewDisk(uint64(anyDiskSize))

	t.Run("ReturnNoErrorIfUserHasNoDisk", func(t *testing.T) {
		err := user.AddDisk(anyDisk)

		assertNoError(t, err)
	})

	t.Run("ReturnErrUserAlreadyHasADiskIfHasDisk", func(t *testing.T) {
		_ = user.AddDisk(anyDisk)

		err := user.AddDisk(anyDisk)

		assertError(t, err)
	})
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error but there is one")
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("Expected an error but there was none")
	}
}
