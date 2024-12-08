package application_test

import (
	"reflect"
	"testing"

	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/mocks"
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

func TestCreateDisk(t *testing.T) {
	anySizeInMib := 1024
	anyUserEmail := "EMAIL@TEST.com"
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	userInRepository := models.NewUser(userInRepoEmail, anyUserPassword)
	repoWithoutUserMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return nil
	}}

	repoWithUserMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return userInRepository
	}}

	t.Run("ReturnErrUserDoesNotExist", func(t *testing.T) {
		service := application.NewCreateDiskService(&repoWithoutUserMock, uint64(anySizeInMib))

		err := service.CreateDisk(anyUserEmail)

		assertError(t, err)
	})

	t.Run("AddNewDiskWithCorrectSize", func(t *testing.T) {
		specifiedDiskSize := 2048
		service := application.NewCreateDiskService(&repoWithUserMock, uint64(specifiedDiskSize))

		service.CreateDisk(anyUserEmail)

		createdDiskSize := userInRepository.GetDiskSpaceLeft()
		assertEquals(t, createdDiskSize, uint64(specifiedDiskSize))
	})

	t.Run("ReturnAddDiskErrValue", func(t *testing.T) {
		service := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMib))
		d := models.NewDisk(uint64(anySizeInMib))
		userInRepository.AddDisk(d)

		err := service.CreateDisk(anyUserEmail)

		assertError(t, err)
	})
}

func assertEquals(t *testing.T, got uint64, want uint64) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got response %d. Should've been %d", got, want)
	}
}
