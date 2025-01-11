package application

import (
	"github.com/Joey-Boivin/sdisk/internal/models"
	"github.com/Joey-Boivin/sdisk/internal/ports"
)

type FetchUserService struct {
	userRepository ports.UserRepository
}

func NewFetchUserService(userRepository ports.UserRepository) *FetchUserService {
	return &FetchUserService{
		userRepository: userRepository,
	}
}

func (fetchUserService *FetchUserService) FetchUser(id string) (models.User, error) {
	userID, err := models.FromString(id)

	if err != nil {
		return models.User{}, err
	}

	user := fetchUserService.userRepository.GetByID(userID)
	if user == nil {
		return models.User{}, &ErrUserDoesNotExist{}
	}

	return *user, nil
}
