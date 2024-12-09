package application_test

import (
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/mocks"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func TestRegisterService(t *testing.T) {
	userInRepoEmail := "John_doe@test.com"
	anyUserEmail := "EMAIL@TEST.com"
	anyUserPassword := "12345"

	userRepoSpy := mocks.UserRepositoryMock{}
	userInRepoMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return models.NewUser(userInRepoEmail, anyUserPassword)
	}}
	registerService := application.NewRegisterService(&userRepoSpy)

	t.Run("UseRegisterServiceToSaveUser", func(t *testing.T) {
		err := registerService.RegisterUser(anyUserEmail, anyUserPassword)

		assertTrue(t, userRepoSpy.SaveUserCalled)
		assertNoError(t, err)
	})

	t.Run("CreateUserWithCorrectParameters", func(t *testing.T) {
		err := registerService.RegisterUser(anyUserEmail, anyUserPassword)

		emailUsed := userRepoSpy.SaveUserCalledWith.GetEmail()
		passwordUsed := userRepoSpy.SaveUserCalledWith.GetPassword()
		assertStringEquals(t, anyUserEmail, emailUsed)
		assertStringEquals(t, anyUserPassword, passwordUsed)
		assertNoError(t, err)
	})

	t.Run("ReturnErrorIfUserAlreadyExists", func(t *testing.T) {
		registerService = application.NewRegisterService(&userInRepoMock)

		err := registerService.RegisterUser(userInRepoEmail, anyUserPassword)

		assertError(t, err)
	})
}

func assertTrue(t *testing.T, statement bool) {
	t.Helper()

	if !statement {
		t.Fatalf("Expected true, got false.")
	}
}

func assertStringEquals(t *testing.T, want string, got string) {
	t.Helper()

	if want != got {
		t.Fatalf("Expected %s, got %s", want, got)
	}
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
