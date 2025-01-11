package application

import (
	"github.com/Joey-Boivin/sdisk/internal/models"
	"github.com/Joey-Boivin/sdisk/internal/ports"
)

type CreateDiskService struct {
	userRepository ports.UserRepository
	sizeInMiB      uint64
	realTimeServer ports.RealTimeServer
}

func NewCreateDiskService(userRepository ports.UserRepository, sizeInMib uint64, server ports.RealTimeServer) *CreateDiskService {
	return &CreateDiskService{
		userRepository: userRepository,
		sizeInMiB:      sizeInMib,
		realTimeServer: server,
	}
}

func (c *CreateDiskService) CreateDisk(id string) error {
	userID, err := models.FromString(id)

	if err != nil {
		return err
	}

	u := c.userRepository.GetByID(userID)
	if u == nil {
		return &ErrUserDoesNotExist{}
	}

	d := models.NewDisk(c.sizeInMiB)
	err = u.AddDisk(d)
	if err != nil {
		return err
	}

	return c.realTimeServer.PrepareDisk(d)
}
