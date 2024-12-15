package infrastructure

import (
	"github.com/Joey-Boivin/sdisk/internal/models"
)

type RamRepository struct {
	users map[string]*models.User
}

func NewRamRepository() *RamRepository {
	repo := RamRepository{}
	repo.users = make(map[string]*models.User)
	return &repo
}

func (r *RamRepository) SaveUser(u *models.User) {
	r.users[u.GetEmail()] = u
}

func (r *RamRepository) GetUser(id string) *models.User {
	if val, ok := r.users[id]; ok {
		return val
	} else {
		return nil
	}
}
