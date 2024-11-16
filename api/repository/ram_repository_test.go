package repository_test

import (
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"github.com/Joey-Boivin/sdisk-api/api/repository"
	"testing"
)

func TestRamRepository(t *testing.T) {

	ramRepository := repository.NewRamRepository()

	t.Run("SaveCorrectUser", func(t *testing.T) {
		anyUserEmail := "EMAIL@TEST.com"
		anyUserPassword := "12345"
		userToSave := models.NewUser(anyUserEmail, anyUserPassword)

		ramRepository.SaveUser(userToSave)

		userSaved := ramRepository.GetUser(anyUserEmail)
		assertSameUsers(t, userToSave, userSaved)
	})

	t.Run("ReturnNilIfUserDoesNotExist", func(t *testing.T) {
		emailOfUserNotInRepository := "JOHN_DOE@TEST.com"

		retrievedUser := ramRepository.GetUser(emailOfUserNotInRepository)

		assertUserIsNil(t, retrievedUser)
	})
}

func assertSameUsers(t *testing.T, got *models.User, want *models.User) {
	if got.GetEmail() != want.GetEmail() {
		t.Fatalf("Expected %p, got %p.", got, want)
	}
}

func assertUserIsNil(t *testing.T, got *models.User) {
	if got != nil {
		t.Fatalf("Expected nil value, got %p.", got)
	}
}
