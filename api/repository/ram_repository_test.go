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
		assertSameUsers(t, userToSave, &userSaved)
	})
}

func assertSameUsers(t *testing.T, got *models.User, want *models.User) {
	if got.GetEmail() != want.GetEmail() {
		t.Fatalf("Expected %p, got %p.", got, want)
	}
}
