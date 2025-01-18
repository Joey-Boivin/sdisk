package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/handlers"
	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/mocks"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

var userInRepoEmail = "John_doe@test.com"
var anyUserPassword = "12345"
var anySizeInMiB = uint64(1024)

var userInRepository = &models.User{}
var idOfUserInRepository = models.UserID{}

var registerService = &application.RegisterService{}
var fetchUserService = &application.FetchUserService{}
var createDiskService = &application.CreateDiskService{}

var userHandler = &handlers.UserHandler{}

var userRepoEmptyMock = mocks.UserRepositoryMock{}
var userRepoWithUserMock = mocks.UserRepositoryMock{}
var serverDummy = mocks.ServerMock{}
var serverMockThatFails = mocks.ServerMock{}

func setup() {
	userInRepository = models.NewUser(userInRepoEmail, anyUserPassword)
	idOfUserInRepository = userInRepository.GetID()

	userRepoEmptyMock = mocks.UserRepositoryMock{FnGetUserByID: func(id models.UserID) *models.User {
		return nil
	}, FnGetUserByEmail: func(email string) *models.User { return nil }}
	userRepoWithUserMock = mocks.UserRepositoryMock{FnGetUserByID: func(id models.UserID) *models.User {
		return userInRepository
	}, FnGetUserByEmail: func(email string) *models.User { return userInRepository }}

	serverDummy = mocks.ServerMock{FnPrepareDisk: func(d *models.Disk) error {
		return nil
	}}

	serverMockThatFails = mocks.ServerMock{FnPrepareDisk: func(d *models.Disk) error {
		anyOpcode := 12
		return &infrastructure.ErrUnknownPacket{Opcode: uint8(anyOpcode)}
	}}

	registerService = application.NewRegisterService(&userRepoWithUserMock)
	fetchUserService = application.NewFetchUserService(&userRepoWithUserMock)
	createDiskService = application.NewCreateDiskService(&userRepoWithUserMock, anySizeInMiB, &serverDummy)
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

	idOfUserInRepository := userInRepository.GetID()

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
		getRequest.SetPathValue("id", idOfUserInRepository.ToString())

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
		getRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.GetUserResource(response, getRequest)

		got := response.Body
		want := createUserJson(userInRepoEmail)
		assertExpectedJson(t, *got, want)
	})

	t.Run("DontIncludeDiskSizeIfUserHasNoDisk", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)
		getRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.GetUserResource(response, getRequest)

		assertResponseHasNoDiskSize(t, response)
	})

	t.Run("IncludeDiskSizeIfUserHasDisk", func(t *testing.T) {
		setup()
		anyDisk := models.NewDisk(anySizeInMiB)
		_ = userInRepository.AddDisk(anyDisk)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.GetUserEndpoint, reader)
		getRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.GetUserResource(response, getRequest)

		assertResponseHasDiskSize(t, response)
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
		postRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusForbidden
		assertStatus(t, got, want)
	})

	t.Run("IfUserDoesNotExistReturnHttpNotFound", func(t *testing.T) {
		setup()
		createDiskService := application.NewCreateDiskService(&userRepoEmptyMock, anySizeInMiB, &serverDummy)
		userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("IfDiskCreatedSuccessReturnHttpStatusCreated", func(t *testing.T) {
		setup()
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusCreated
		assertStatus(t, got, want)
	})

	t.Run("IfServerFailsReturnInternalError", func(t *testing.T) {
		setup()
		createDiskService = application.NewCreateDiskService(&userRepoWithUserMock, anySizeInMiB, &serverMockThatFails)
		userHandler := handlers.NewUserHandler(registerService, fetchUserService, createDiskService)
		response := httptest.NewRecorder()
		reader := strings.NewReader("")
		postRequest, _ := http.NewRequest(http.MethodPost, handlers.CreateDiskEndpoint, reader)
		postRequest.SetPathValue("id", idOfUserInRepository.ToString())

		userHandler.CreateDiskResource(response, postRequest)

		got := response.Code
		want := http.StatusInternalServerError
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
	t.Helper()

	if reflect.DeepEqual(got.Bytes(), want.Bytes()) {
		t.Fatalf("Expected the following json:\n%s\nBut got:\n%s", got.String(), want.String())
	}
}

func responseHasDiskSize(t *testing.T, response *httptest.ResponseRecorder) bool {
	t.Helper()

	var responseBody = make(map[string]interface{})
	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	_, exists := responseBody["diskSpaceInMiB"]

	return exists
}

func assertResponseHasNoDiskSize(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if responseHasDiskSize(t, response) {
		t.Fatalf("Expected no disk size field but there was one")
	}
}

func assertResponseHasDiskSize(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if !responseHasDiskSize(t, response) {
		t.Fatalf("Expected disk size field but there was none")
	}
}

func createUserJson(email string) bytes.Buffer {
	var buf bytes.Buffer
	str := fmt.Sprintf("{\"email\": \"%s\"}", email)
	buf.WriteString(str)
	return buf
}
