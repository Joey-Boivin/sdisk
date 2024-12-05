package handlers

import (
	"encoding/json"
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"net/http"
)

const (
	CreateUserEndpoint = "POST /users"
	GetUserEndpoint    = "GET /users/{id}"
)

type UserHandler struct {
	registerService  *application.RegisterService
	fetchUserService *application.FetchUserService
}

type RegisterRequest struct {
	Email    string
	Password string
}

type FetchUserResponse struct {
	Email string `json:"email"`
}

func NewUserHandler(registerService *application.RegisterService, fetchUserService *application.FetchUserService) *UserHandler {
	return &UserHandler{
		registerService:  registerService,
		fetchUserService: fetchUserService,
	}
}

func (r *UserHandler) Post(writer http.ResponseWriter, req *http.Request) {
	var registerRequest RegisterRequest

	err := json.NewDecoder(req.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.registerService.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (r *UserHandler) Get(writer http.ResponseWriter, req *http.Request) {
	email := req.PathValue("id")
	user, err := r.fetchUserService.FetchUser(email)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	resp := FetchUserResponse{Email: user.GetEmail()}

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
