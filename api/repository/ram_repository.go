package repository

import "github.com/Joey-Boivin/sdisk-api/api/models"

type RamRepository struct {
	users map[string]models.User
}

func NewRamRepository() *RamRepository {
	repo := RamRepository{}
	repo.users = make(map[string]models.User)
	return &repo
}

func (r *RamRepository) SaveUser(u *models.User) {
	r.users[u.GetEmail()] = *u
}

func (r *RamRepository) GetUser(id string) models.User {
	return r.users[id]
}
