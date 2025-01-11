package application_test

import (
	"errors"
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/mocks"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func TestCreateDisk(t *testing.T) {
	anySizeInMiB := uint64(1024)
	anyUserEmail := "EMAIL@TEST.com"
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	userInRepository := models.NewUser(userInRepoEmail, anyUserPassword)
	idOfUserInRepository := userInRepository.GetID()

	repoWithoutUserMock := mocks.UserRepositoryMock{FnGetUserByID: func(id models.UserID) *models.User {
		return nil
	}}

	repoWithUserMock := mocks.UserRepositoryMock{FnGetUserByID: func(id models.UserID) *models.User {
		return userInRepository
	}}

	serverMockThatFails := mocks.ServerMock{FnPrepareDisk: func(d *models.Disk) error {
		return errors.New("server failed to prepare disk")
	}}

	serverMockDummy := mocks.ServerMock{FnPrepareDisk: func(d *models.Disk) error {
		return nil
	}}

	t.Run("ReturnErrUserDoesNotExist", func(t *testing.T) {
		service := application.NewCreateDiskService(&repoWithoutUserMock, anySizeInMiB, &serverMockDummy)

		err := service.CreateDisk(anyUserEmail)

		assertError(t, err)
	})

	t.Run("AddNewDiskWithCorrectSize", func(t *testing.T) {
		specifiedDiskSize := uint64(2048)
		service := application.NewCreateDiskService(&repoWithUserMock, specifiedDiskSize, &serverMockDummy)

		_ = service.CreateDisk(idOfUserInRepository.ToString())

		createdDiskSize, _ := userInRepository.GetDiskSpaceLeft()
		assertEquals(t, createdDiskSize, uint64(specifiedDiskSize))
	})

	t.Run("ReturnAddDiskErrValue", func(t *testing.T) {
		service := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMiB), &serverMockDummy)
		d := models.NewDisk(uint64(anySizeInMiB))
		_ = userInRepository.AddDisk(d)

		err := service.CreateDisk(idOfUserInRepository.ToString())

		assertError(t, err)
	})

	t.Run("ReturnServerFailureError", func(t *testing.T) {
		userInRepository = models.NewUser(userInRepoEmail, anyUserPassword)
		service := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMiB), &serverMockThatFails)

		err := service.CreateDisk(idOfUserInRepository.ToString())

		assertError(t, err)
	})

	t.Run("ServerPreparesDisk", func(t *testing.T) {
		userInRepository = models.NewUser(userInRepoEmail, anyUserPassword)
		service := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMiB), &serverMockDummy)

		_ = service.CreateDisk(idOfUserInRepository.ToString())

		assertTrue(t, serverMockDummy.PrepareDiskCalled)
	})
}

func assertEquals(t *testing.T, got uint64, want uint64) {
	t.Helper()

	if got != want {
		t.Fatalf("Got response %d. Should've been %d", got, want)
	}
}
