package infrastructure_test

import (
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

func TestSaveUser(t *testing.T) {
	ramRepository := infrastructure.NewRamRepository()
	anyUserEmail := "EMAIL@TEST.com"
	var userToSave *models.User

	t.Run("SaveCorrectUser", func(t *testing.T) {
		anyUserPassword := "12345"
		userToSave = models.NewUser(anyUserEmail, anyUserPassword)
		anyUserID := userToSave.GetID()

		ramRepository.SaveUser(userToSave)

		userSaved := ramRepository.GetByID(anyUserID)
		assertSameUsers(t, userToSave, userSaved)
	})

	t.Run("UserCanBeRetrievedByEmail", func(t *testing.T) {
		anyUserPassword := "12345"
		userToSave = models.NewUser(anyUserEmail, anyUserPassword)

		ramRepository.GetByEmail(anyUserEmail)

		userSaved := ramRepository.GetByEmail(anyUserEmail)
		assertSameUsers(t, userToSave, userSaved)
	})
}

func TestGetUser(t *testing.T) {
	ramRepository := infrastructure.NewRamRepository()

	t.Run("ReturnNilIfUserDoesNotExist", func(t *testing.T) {
		idOfAUserNotInRepository := models.NewUserID()

		retrievedUser := ramRepository.GetByID(idOfAUserNotInRepository)

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
