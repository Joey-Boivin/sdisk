package application_test

import (
	"testing"
)

func TestNewFetchUserService(t *testing.T) {
	/*
		userInRepoEmail := "John_doe@test.com"
		anyUserEmail := "EMAIL@TEST.com"
		anyUserPassword := "12345"

		userRepoDummy := mocks.RamRepository{}
		userInRepoMock := mocks.RamRepository{FnGetUser: func(id string) *models.User {
			return models.NewUser(userInRepoEmail, anyUserPassword)
		}}
		fetchUserService := application.NewFetchUserService(&userRepoDummy)

		t.Run("ReturnErrorIfUserDoesNotExist", func(t *testing.T) {
			_, err := fetchUserService.FetchUser(anyUserEmail)

			assertError(t, err)
		})

		t.Run("ReturnNoErrorIfUserExists", func(t *testing.T) {
			fetchUserService = application.NewFetchUserService(&userInRepoMock)

			_, err := fetchUserService.FetchUser(anyUserEmail)

			assertNoError(t, err)
		})

		t.Run("ReturnUserFromRepoIfItExists", func(t *testing.T) {
			fetchUserService = application.NewFetchUserService(&userInRepoMock)

			user, _ := fetchUserService.FetchUser(userInRepoEmail)

			assertStringEquals(t, userInRepoEmail, user.GetEmail())
		})
	*/
}
