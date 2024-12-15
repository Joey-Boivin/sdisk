package ports

import "github.com/Joey-Boivin/sdisk/internal/models"

type UserRepository interface {
	SaveUser(u *models.User)
	GetUser(id string) *models.User
}

type RealTimeServer interface {
	PrepareDisk(d *models.Disk) error
}
