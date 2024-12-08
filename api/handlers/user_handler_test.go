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

func TestCreateUser(t *testing.T) {
	validUserJson := "{\"email\": \"EMAIL@TEST.com\", \"password\": \"12345\"}"
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	anySizeInMib := 1024
	userRepoDummy := mocks.UserRepositoryMock{}
	registerService := application.NewRegisterService(&userRepoDummy)
	fetchUserService := application.NewFetchUserService(&userRepoDummy)
	createDiskService := application.NewCreateDiskService(&userRepoDummy, uint64(anySizeInMib))
	userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
	userInRepoMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return models.NewUser(userInRepoEmail, anyUserPassword)
	}}

	t.Run("ReturnHttpBadRequestIfParseError", func(t *testing.T) {
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
		response := httptest.NewRecorder()
		badRequest := "{"
		reader := strings.NewReader(badRequest)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertNoSave(t, userRepoDummy)
	})

	t.Run("ReturnHttpCreatedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

	t.Run("UserSavedIfNoErrors", func(t *testing.T) {
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertSave(t, userRepoDummy)
	})

	t.Run("ReturnHttpForbiddenIfUserAlreadyExists", func(t *testing.T) {
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()

		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("DoNotOverrideExistingUser", func(t *testing.T) {
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader(validUserJson)
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateUserEndpoint, reader)

		userHandler.CreateUserResource(response, postRequest)

		assertNoSave(t, userInRepoMock)
	})
}

func TestGetUser(t *testing.T) {
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	anySizeInMib := 1024
	userRepoDummy := mocks.UserRepositoryMock{}
	registerService := application.NewRegisterService(&userRepoDummy)
	fetchUserService := application.NewFetchUserService(&userRepoDummy)
	createDiskService := application.NewCreateDiskService(&userRepoDummy, uint64(anySizeInMib))
	userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
	userInRepoMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
		return models.NewUser(userInRepoEmail, anyUserPassword)
	}}

	t.Run("ReturnHttpNotFoundIfUserDoesNotExist", func(t *testing.T) {
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.CreateUserEndpoint, reader)

		userHandler.GetUserResource(response, getRequest)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("ReturnHttpOkIfUserExists", func(t *testing.T) {
		fetchUserService = application.NewFetchUserService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
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
		registerService = application.NewRegisterService(&userInRepoMock)
		userHandler = handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
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
	userInRepoEmail := "John_doe@test.com"
	anyUserPassword := "12345"
	anySizeInMib := 1024
	dummyUserRepository := mocks.UserRepositoryMock{}
	t.Run("IfUserHasADiskReturnHttpStatusForbidden", func(t *testing.T) {
		userInRepository := models.NewUser(userInRepoEmail, anyUserPassword)
		repoWithUserMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
			return userInRepository
		}}
		existingDisk := models.NewDisk(uint64(anySizeInMib))
		userInRepository.AddDisk(existingDisk)
		registerService := application.NewRegisterService(&repoWithUserMock)
		fetchUserService := application.NewFetchUserService(&dummyUserRepository)
		createDiskService := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMib))
		userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
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
		registerService := application.NewRegisterService(&dummyUserRepository)
		fetchUserService := application.NewFetchUserService(&dummyUserRepository)
		createDiskService := application.NewCreateDiskService(&dummyUserRepository, uint64(anySizeInMib))
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
		userInRepository := models.NewUser(userInRepoEmail, anyUserPassword)
		repoWithUserMock := mocks.UserRepositoryMock{FnGetUser: func(id string) *models.User {
			return userInRepository
		}}
		registerService := application.NewRegisterService(&dummyUserRepository)
		fetchUserService := application.NewFetchUserService(&dummyUserRepository)
		createDiskService := application.NewCreateDiskService(&repoWithUserMock, uint64(anySizeInMib))
		userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
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
