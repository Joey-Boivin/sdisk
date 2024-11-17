package mocks

import "github.com/Joey-Boivin/sdisk-api/api/models"

type RamRepository struct {
	FnSaveUser         func(u *models.User)
	SaveUserCalled     bool
	SaveUserCalledWith *models.User

	FnGetUser         func(id string) *models.User
	GetUserCalled     bool
	GetUserCalledWith string
}

func (r *RamRepository) SaveUser(u *models.User) {
	r.SaveUserCalled = true
	r.SaveUserCalledWith = u

	if r.FnSaveUser != nil {
		r.FnSaveUser(u)
	}
}

func (r *RamRepository) GetUser(id string) *models.User {
	r.GetUserCalled = true
	r.GetUserCalledWith = id

	if r.FnGetUser != nil {
		return r.FnGetUser(id)
	}
	return nil
}
