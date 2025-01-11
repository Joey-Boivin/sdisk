package application

import (
	"github.com/Joey-Boivin/sdisk/internal/models"
	"github.com/Joey-Boivin/sdisk/internal/ports"
)

type RegisterService struct {
	userRepository ports.UserRepository
}

func NewRegisterService(userRepository ports.UserRepository) *RegisterService {
	return &RegisterService{
		userRepository: userRepository,
	}
}

func (registerService *RegisterService) RegisterUser(email string, password string) (models.UserID, error) {
	user := registerService.userRepository.GetByEmail(email)

	if user != nil {
		return user.GetID(), &ErrUserAlreadyExists{email}
	}

	user = models.NewUser(email, password)
	registerService.userRepository.SaveUser(user)

	return user.GetID(), nil
}
