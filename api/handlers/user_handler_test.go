package handlers_test

import (
	"bytes"
	"fmt"
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/handlers"
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"github.com/Joey-Boivin/sdisk-api/api/repository/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	userRepoDummy := mocks.RamRepository{}
	registerService := application.NewRegisterService(&userRepoDummy)
	fetchUserService := application.NewFetchUserService(&userRepoDummy)
	userHandler := handlers.NewUserHandler(registerService, fetchUserService)

	t.Run("ReturnHttpBadRequestIfParseError", func(t *testing.T) {
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusBadRequest
		assertStatus(t, got, want)
	})

	t.Run("DoNotSaveAUserIfParseError", func(t *testing.T) {
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertNoSave(t, userRepoDummy)
	})

	t.Run("ReturnHttpCreatedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

	t.Run("UserSavedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertSave(t, userRepoDummy)
	})

	t.Run("ReturnHttpForbiddenIfUserAlreadyExists", func(t *testing.T) {
		userInRepoMock := mocks.RamRepository{FnGetUser: func(id string) *models.User {
			return models.NewUser(userInRepoEmail, anyUserPassword)
		}}
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService)
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("DoNotOverrideExistingUser", func(t *testing.T) {
		userInRepoMock := mocks.RamRepository{FnGetUser: func(id string) *models.User {
			return models.NewUser(userInRepoEmail, anyUserPassword)
		}}
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService)
		response := httptest.NewRecorder()
		validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.Post(response, postRequest)

		assertNoSave(t, userInRepoMock)
	})
}

func GetUser(t *testing.T) {
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	userRepoDummy := mocks.RamRepository{}
	registerService := application.NewRegisterService(&userRepoDummy)
	fetchUserService := application.NewFetchUserService(&userRepoDummy)
	userHandler := handlers.NewUserHandler(registerService, fetchUserService)

	t.Run("ReturnHttpNotFoundIfUserDoesNotExist", func(t *testing.T) {
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.CreateUserEndpoint, reader)

		userHandler.Get(response, getRequest)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("ReturnHttpOkIfUserExists", func(t *testing.T) {
		userInRepoMock := mocks.RamRepository{FnGetUser: func(id string) *models.User {
			return models.NewUser(userInRepoEmail, anyUserPassword)
		}}
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)

		userHandler.Get(response, getRequest)

		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})

	t.Run("ReturnExpectedJsonIfUserExists", func(t *testing.T) {
		userInRepoMock := mocks.RamRepository{FnGetUser: func(id string) *models.User {
			return models.NewUser(userInRepoEmail, anyUserPassword)
		}}
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)

		userHandler.Get(response, getRequest)

		got := response.Body
		want := createUserJson(userInRepoEmail)
		assertExpectedJson(t, *got, want)
	})
}

func assertNoSave(t *testing.T, userRepoMock mocks.RamRepository) {
	t.Helper()

	if userRepoMock.SaveUserCalled {
		t.Fatalf("No save should have happened")
	}
}

func assertSave(t *testing.T, userRepoMock mocks.RamRepository) {
	t.Helper()

	if !userRepoMock.SaveUserCalled {
		t.Fatalf("User should have been saved but was not")
	}
}

func assertExpectedJson(t *testing.T, got bytes.Buffer, want bytes.Buffer) {
	if reflect.DeepEqual(got.Bytes(), want.Bytes()) {
		t.Fatalf("Expected the following json:\n%s\nBut got:\n%s", got.String(), want.String())
	}
}

func createUserJson(email string) bytes.Buffer {
	var buf bytes.Buffer
	str := fmt.Sprintf("{\"email\": \"%s\"}", email)
	buf.WriteString(str)
	return buf
}
