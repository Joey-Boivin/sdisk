package application_test

import (
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/mocks"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func TestNewFetchUserService(t *testing.T) {
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"

	userRepoDummy := mocks.UserRepositoryMock{}
	userInRepo := models.NewUser(userInRepoEmail, anyUserPassword)
	idOfUserInRepository := userInRepo.GetID()

	userInRepoMock := mocks.UserRepositoryMock{FnGetUserByID: func(id models.UserID) *models.User {
		return userInRepo
	}}
	fetchUserService := application.NewFetchUserService(&userRepoDummy)

	t.Run("ReturnErrorIfUserDoesNotExist", func(t *testing.T) {
		_, err := fetchUserService.FetchUser(idOfUserInRepository.ToString())

		assertError(t, err)
	})

	t.Run("ReturnNoErrorIfUserExists", func(t *testing.T) {
		fetchUserService = application.NewFetchUserService(&userInRepoMock)

		_, err := fetchUserService.FetchUser(idOfUserInRepository.ToString())

		assertNoError(t, err)
	})

	t.Run("ReturnUserFromRepoIfItExists", func(t *testing.T) {
		fetchUserService = application.NewFetchUserService(&userInRepoMock)

		user, _ := fetchUserService.FetchUser(idOfUserInRepository.ToString())

		assertStringEquals(t, userInRepoEmail, user.GetEmail())
	})
}
