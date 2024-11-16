package application_test

import (
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"testing"
)

type userRepositoryMock struct {
	used         bool
	emailUsed    string
	passwordUsed string
}

func (r *userRepositoryMock) SaveUser(u *models.User) {
	r.used = true
	r.emailUsed = u.GetEmail()
	r.passwordUsed = u.GetPassword()
}

func (r *userRepositoryMock) GetUser(id string) models.User {
	return *models.NewUser("1", "1")
}

func TestRegisterService(t *testing.T) {
	userRepoMock := userRepositoryMock{used: false}
	registerService := application.NewRegisterService(&userRepoMock)

	t.Run("CreateUserWithCorrectParameters", func(t *testing.T) {
		userEmail := "EMAIL@TEST.com"
		userPassword := "12345"

		registerService.RegisterUser(userEmail, userPassword)

		assertStringEquals(t, userEmail, userRepoMock.emailUsed)
		assertStringEquals(t, userPassword, userRepoMock.passwordUsed)
	})

	t.Run("UseRegisterServiceToSaveUser", func(t *testing.T) {
		anyUserEmail := "EMAIL@TEST.com"
		anyUserPassword := "12345"

		registerService.RegisterUser(anyUserEmail, anyUserPassword)

		assertTrue(t, userRepoMock.used)
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
