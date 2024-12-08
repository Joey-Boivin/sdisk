package application

import (
	"github.com/Joey-Boivin/sdisk-api/api/models"
	"github.com/Joey-Boivin/sdisk-api/api/ports"
)

type FetchUserService struct {
	userRepository ports.UserRepository
}

func NewFetchUserService(userRepository ports.UserRepository) *FetchUserService {
	return &FetchUserService{
		userRepository: userRepository,
	}
}

func (fetchUserService *FetchUserService) FetchUser(email string) (models.User, error) {
	user := fetchUserService.userRepository.GetUser(email)
	if user == nil {
		return models.User{}, &ErrUserDoesNotExist{}
	}

	return *user, nil
}
