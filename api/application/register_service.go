package application

import (
	"fmt"
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

func CreateErrorUserAlreadyExists(email string) error {
	return fmt.Errorf("user already exists for the email address %s", email)
}

type RegisterService struct {
	userRepository models.UserRepository
}

func NewRegisterService(userRepository models.UserRepository) *RegisterService {
	return &RegisterService{
		userRepository: userRepository,
	}
}

func (registerService *RegisterService) RegisterUser(email string, password string) error {
	user := registerService.userRepository.GetUser(email)

	if user != nil {
		return CreateErrorUserAlreadyExists(email)
	}

	user = models.NewUser(email, password)
	registerService.userRepository.SaveUser(user)
	return nil
}
