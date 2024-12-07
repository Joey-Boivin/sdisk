package application

import "github.com/Joey-Boivin/sdisk-api/api/models"

type CreateDiskService struct {
	userRepository models.UserRepository
	sizeInMib      uint64
}

func NewCreateDiskService(userRepository models.UserRepository, sizeInMib uint64) *CreateDiskService {
	return &CreateDiskService{
		userRepository: userRepository,
		sizeInMib:      sizeInMib,
	}
}

func (c *CreateDiskService) CreateDisk(email string) error {
	u := c.userRepository.GetUser(email)
	if u == nil {
		return &ErrUserDoesNotExist{}
	}

	d := models.NewDisk(c.sizeInMib)
	err := u.AddDisk(d)
	return err
}
