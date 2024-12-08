package ports

import "github.com/Joey-Boivin/sdisk-api/api/models"

type UserRepository interface {
	SaveUser(u *models.User)
	GetUser(id string) *models.User
}
