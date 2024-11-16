package handlers

import (
	"encoding/json"
	"github.com/Joey-Boivin/sdisk-api/api/application"
	"net/http"
)

const (
	UsersEndpoint = "/users"
)

type UserHandler struct {
	service *application.RegisterService
}

type RegisterRequest struct {
	Email    string
	Password string
}

func NewUserHandler(service *application.RegisterService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (r *UserHandler) Post(writer http.ResponseWriter, req *http.Request) {
	var registerRequest RegisterRequest

	err := json.NewDecoder(req.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.service.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}
