package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Joey-Boivin/sdisk/internal/application"
	"github.com/Joey-Boivin/sdisk/internal/infrastructure"
	"github.com/Joey-Boivin/sdisk/internal/models"
)

const (
	CreateUserEndpoint = "POST /users"
	GetUserEndpoint    = "GET /users/{id}"
	CreateDiskEndpoint = "POST /users/{id}/disk"
)

type UserHandler struct {
	registerService   *application.RegisterService
	fetchUserService  *application.FetchUserService
	createDiskService *application.CreateDiskService
}

type RegisterRequest struct {
	Email    string
	Password string
}

type FetchUserResponse struct {
	Email     string `json:"email"`
	DiskSpace int    `json:"diskSpaceInMiB,omitempty"`
}

func NewUserHandler(registerService *application.RegisterService, fetchUserService *application.FetchUserService, createDiskService *application.CreateDiskService) *UserHandler {
	return &UserHandler{
		registerService:   registerService,
		fetchUserService:  fetchUserService,
		createDiskService: createDiskService,
	}
}

func (h *UserHandler) CreateUserResource(writer http.ResponseWriter, req *http.Request) {
	var registerRequest RegisterRequest

	err := json.NewDecoder(req.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.registerService.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writer.Header().Set("Location", fmt.Sprintf("/users/%s", id.ToString()))
	writer.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUserResource(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	user, err := h.fetchUserService.FetchUser(id)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	space, err := user.GetDiskSpaceLeft()
	var resp FetchUserResponse

	if err != nil {
		resp = FetchUserResponse{Email: user.GetEmail()}
	} else {
		resp = FetchUserResponse{Email: user.GetEmail(), DiskSpace: int(space)}
	}

	data, err := json.Marshal(resp)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
}

func (h *UserHandler) CreateDiskResource(writer http.ResponseWriter, req *http.Request) {
	email := req.PathValue("id")

	err := h.createDiskService.CreateDisk(email)

	if err != nil {
		switch err.(type) {
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			return

		case *models.ErrUserAlreadyHasADisk:
			writer.WriteHeader(http.StatusForbidden)
			return

		case *application.ErrUserDoesNotExist:
			writer.WriteHeader(http.StatusNotFound)
			return

		case *infrastructure.ErrUnknownJob:
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	writer.WriteHeader(http.StatusCreated)
}
