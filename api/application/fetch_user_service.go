package application

import (
	"fmt"
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

func CreateErrorUserDoesNotExist(email string) error {
	return fmt.Errorf("user with email %s does not exist", email)
}

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
		return models.User{}, CreateErrorUserDoesNotExist(email)
	}

	return *user, nil
}
