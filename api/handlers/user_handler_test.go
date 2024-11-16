package handlers_test

import (
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/handlers"
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestUserHandler(t *testing.T) {
	userRepoMock := userRepositoryMock{}
	service := application.NewRegisterService(&userRepoMock)
	userHandler := handlers.NewUserHandler(service)

	t.Run("ReturnHttpBadRequestIfParseError", func(t *testing.T) {
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusBadRequest
		assertStatus(t, got, want)
	})

	t.Run("ReturnHttpOkIfNoParseErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

}
