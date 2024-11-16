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
	userAlreadyExists bool
	saved             bool
}

func (r *userRepositoryMock) SaveUser(u *models.User) {
	r.saved = true
}

func (r *userRepositoryMock) GetUser(id string) *models.User {
	if r.userAlreadyExists {
		anyUserPassword := "1"
		return models.NewUser(id, anyUserPassword)
	}

	return nil
}

func TestBadRequest(t *testing.T) {
	userRepoMock := userRepositoryMock{saved: false}
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

	t.Run("DoNotSaveAUserIfParseError", func(t *testing.T) {
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertNoSave(t, userRepoMock)
	})
}

func TestUserAlreadyExists(t *testing.T) {
	userRepoMock := userRepositoryMock{saved: false, userAlreadyExists: true}
	service := application.NewRegisterService(&userRepoMock)
	userHandler := handlers.NewUserHandler(service)

	t.Run("ReturnHttpForbiddenIfUserAlreadyExists", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("DoNotOverrideExistingUser", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertNoSave(t, userRepoMock)
	})
}

func TestUserCanRegister(t *testing.T) {
	userRepoMock := userRepositoryMock{}
	service := application.NewRegisterService(&userRepoMock)
	userHandler := handlers.NewUserHandler(service)

	t.Run("ReturnHttpCreatedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

	t.Run("UserSavedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.UsersEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertSave(t, userRepoMock)
	})
}

func assertNoSave(t *testing.T, userRepoMock userRepositoryMock) {
	if userRepoMock.saved {
		t.Fatalf("No save should have happened")
	}
}

func assertSave(t *testing.T, userRepoMock userRepositoryMock) {
	if !userRepoMock.saved {
		t.Fatalf("User should have been saved but was not")
	}
}
