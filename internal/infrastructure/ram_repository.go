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
	id := u.GetID()
	r.users[id.ToString()] = u
}

func (r *RamRepository) GetByID(id models.UserID) *models.User {
	if val, ok := r.users[id.ToString()]; ok {
		return val
	} else {
		return nil
	}
}

func (r *RamRepository) GetByEmail(email string) *models.User {
	for _, val := range r.users {
		if val.GetEmail() == email {
			return val
		}
	}

	return nil
}
