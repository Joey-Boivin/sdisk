package application_test

import (
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"testing"
)

type userRepositoryMock struct {
	used              bool
	emailUsed         string
	passwordUsed      string
	userAlreadyExists bool
}

func (r *userRepositoryMock) SaveUser(u *models.User) {
	r.used = true
	r.emailUsed = u.GetEmail()
	r.passwordUsed = u.GetPassword()
}

func (r *userRepositoryMock) GetUser(id string) *models.User {
	if r.userAlreadyExists {
		anyUserPassword := "1"
		return models.NewUser(id, anyUserPassword)
	}

	return nil
}

func TestRegisterService(t *testing.T) {
	userRepoMock := userRepositoryMock{used: false, userAlreadyExists: false}
	registerService := application.NewRegisterService(&userRepoMock)

	t.Run("CreateUserWithCorrectParameters", func(t *testing.T) {
		userEmail := "EMAIL@TEST.com"
		userPassword := "12345"

		err := registerService.RegisterUser(userEmail, userPassword)

		assertStringEquals(t, userEmail, userRepoMock.emailUsed)
		assertStringEquals(t, userPassword, userRepoMock.passwordUsed)
		assertNoError(t, err)
	})

	t.Run("UseRegisterServiceToSaveUser", func(t *testing.T) {
		anyUserEmail := "EMAIL@TEST.com"
		anyUserPassword := "12345"

		err := registerService.RegisterUser(anyUserEmail, anyUserPassword)

		assertTrue(t, userRepoMock.used)
		assertNoError(t, err)
	})

	t.Run("ReturnErrorIfUserAlreadyExists", func(t *testing.T) {
		userInRepoEmail := "JOHN_DOE@TEST.com"
		anyPassword := "123"
		userRepoMock.userAlreadyExists = true

		err := registerService.RegisterUser(userInRepoEmail, anyPassword)

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
