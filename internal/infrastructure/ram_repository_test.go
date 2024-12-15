package infrastructure_test

import (
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func TestSaveUser(t *testing.T) {
	ramRepository := infrastructure.NewRamRepository()

	t.Run("SaveCorrectUser", func(t *testing.T) {
		anyUserEmail := "EMAIL@TEST.com"
		anyUserPassword := "12345"
		userToSave := models.NewUser(anyUserEmail, anyUserPassword)

		ramRepository.SaveUser(userToSave)

		userSaved := ramRepository.GetUser(anyUserEmail)
		assertSameUsers(t, userToSave, userSaved)
	})
}

func TestGetUser(t *testing.T) {
	ramRepository := infrastructure.NewRamRepository()

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
