package application

import (
	"github.com/Joey-Boivin/sdisk/internal/models"
	"github.com/Joey-Boivin/sdisk/internal/ports"
)

type CreateDiskService struct {
	userRepository ports.UserRepository
	sizeInMiB      uint64
}

func NewCreateDiskService(userRepository ports.UserRepository, sizeInMib uint64) *CreateDiskService {
	return &CreateDiskService{
		userRepository: userRepository,
		sizeInMiB:      sizeInMib,
	}
}

func (c *CreateDiskService) CreateDisk(email string) error {
	u := c.userRepository.GetUser(email)
	if u == nil {
		return &ErrUserDoesNotExist{}
	}

	d := models.NewDisk(c.sizeInMiB)
	err := u.AddDisk(d)
	return err
}
