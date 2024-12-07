package application

import (
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

type FetchUserService struct {
	userRepository models.UserRepository
}

func NewFetchUserService(userRepository models.UserRepository) *FetchUserService {
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
