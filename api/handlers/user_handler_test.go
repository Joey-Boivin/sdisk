package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Joey-Boivin/sdisk-api/api/application"
	"github.com/Joey-Boivin/sdisk-api/api/handlers"
	"github.com/Joey-Boivin/sdisk-api/api/mocks"
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

var userInRepoEmail = "John_doe@test.com"
var anyUserPassword = "12345"
var anySizeInMiB = uint64(1024)

var userInRepository = &models.User{}

var registerService = &application.RegisterService{}
var fetchUserService = &application.FetchUserService{}
var createDiskService = &application.CreateDiskService{}

var userHandler = &handlers.UserHandler{}

var userRepoEmptyMock = mocks.UserRepositoryMock{}
var userRepoWithUserMock = mocks.UserRepositoryMock{}

func setup() {
	userInRepository = models.NewUser(userInRepoEmail, anyUserPassword)
	userRepoEmptyMock = mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return nil
	}}
	userRepoWithUserMock = mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return userInRepository
	}}
	registerService = application.NewRegisterService(&userRepoWithUserMock)
	fetchUserService = application.NewFetchUserService(&userRepoWithUserMock)
	createDiskService = application.NewCreateDiskService(&userRepoWithUserMock, anySizeInMiB)
	userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
}

func TestCreateUser(t *testing.T) {
	validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"

	t.Run("ReturnHttpBadRequestIfParseError", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		got := response.Code
		want := http.StatusBadRequest
		assertStatus(t, got, want)
	})

	t.Run("DoNotSaveAUserIfParseError", func(t *testing.T) {
		setup()
		registerService = application.NewRegisterService(&userRepoEmptyMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertNoSave(t, userRepoEmptyMock)
	})

	t.Run("ReturnHttpCreatedIfNoErrors", func(t *testing.T) {
		setup()
		registerService = application.NewRegisterService(&userRepoEmptyMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

	t.Run("UserSavedIfNoErrors", func(t *testing.T) {
		setup()
		registerService = application.NewRegisterService(&userRepoEmptyMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertSave(t, userRepoEmptyMock)
	})

	t.Run("ReturnHttpForbiddenIfUserAlreadyExists", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()

		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("DoNotOverrideExistingUser", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertNoSave(t, userRepoWithUserMock)
	})
}

func TestGetUser(t *testing.T) {

	t.Run("ReturnHttpNotFoundIfUserDoesNotExist", func(t *testing.T) {
		setup()
		fetchUserService = application.NewFetchUserService(&userRepoEmptyMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.CreateUserEndpoint, reader)

		userHandler.GetUserResource(response, getRequest)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("ReturnHttpOkIfUserExists", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)
		getRequest.SetPathValue("id", userInRepoEmail)

		userHandler.GetUserResource(response, getRequest)

		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})

	t.Run("ReturnExpectedJsonIfUserExists", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)

		userHandler.GetUserResource(response, getRequest)

		got := response.Body
		want := createUserJson(userInRepoEmail)
		assertExpectedJson(t, *got, want)
	})
}

func TestCreateDiskResource(t *testing.T) {
	t.Run("IfUserHasADiskReturnHttpStatusForbidden", func(t *testing.T) {
		setup()
		existingDisk := models.NewDisk(uint64(anySizeInMiB))
		_ = userInRepository.AddDisk(existingDisk)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", userInRepoEmail)

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("IfUserDoesNotExistReturnHttpNotFound", func(t *testing.T) {
		setup()
		createDiskService := application.NewCreateDiskService(&userRepoEmptyMock, anySizeInMiB)
		userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", userInRepoEmail)

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("IfUnknownErrorThenReturnHttpStatusInternalServerError", func(t *testing.T) {

	})

	t.Run("IfDiskCreatedSuccessReturnHttpStatusCreated", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", userInRepoEmail)

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

}

func assertNoSave(t *testing.T, userRepoMock mocks.UserRepositoryMock) {
	t.Helper()

	if userRepoMock.SaveUserCalled {
		t.Fatalf("No save should have happened")
	}
}

func assertSave(t *testing.T, userRepoMock mocks.UserRepositoryMock) {
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
