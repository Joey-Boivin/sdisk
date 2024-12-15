package mocks

import "github.com/Joey-Boivin/sdisk/internal/models"

type UserRepositoryMock struct {
	FnSaveUser         func(u *models.User)
	SaveUserCalled     bool
	SaveUserCalledWith *models.User

	FnGetUser         func(id string) *models.User
	GetUserCalled     bool
	GetUserCalledWith string
}

func (r *UserRepositoryMock) SaveUser(u *models.User) {
	r.SaveUserCalled = true
	r.SaveUserCalledWith = u

	if r.FnSaveUser != nil {
		r.FnSaveUser(u)
	}
}

func (r *UserRepositoryMock) GetUser(id string) *models.User {
	r.GetUserCalled = true
	r.GetUserCalledWith = id

	if r.FnGetUser != nil {
		return r.FnGetUser(id)
	}
	return nil
}

type ServerMock struct {
	FnPrepareDisk         func(d *models.Disk) error
	PrepareDiskCalled     bool
	PrepareDiskCalledWith *models.Disk
}

func (s *ServerMock) PrepareDisk(d *models.Disk) error {
	s.PrepareDiskCalledWith = d
	s.PrepareDiskCalled = true

	if s.FnPrepareDisk != nil {
		return s.FnPrepareDisk(d)
	}

	return nil
}
