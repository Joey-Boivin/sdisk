package ports

import "github.com/Joey-Boivin/sdisk/internal/models"

type UserRepository interface {
	SaveUser(u *models.User)
	GetByID(id models.UserID) *models.User
	GetByEmail(email string) *models.User
}

type RealTimeServer interface {
	Run()
	PrepareDisk(d *models.Disk, user *models.User) error
}
