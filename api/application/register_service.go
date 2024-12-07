package application

import (
	"github.com/Joey-Boivin/sdisk-api/api/models"
)

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
		return &ErrUserAlreadyExists{email}
	}

	user = models.NewUser(email, password)
	registerService.userRepository.SaveUser(user)
	return nil
}
