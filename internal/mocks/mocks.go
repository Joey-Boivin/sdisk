package mocks

import "github.com/Joey-Boivin/sdisk/internal/models"

type UserRepositoryMock struct {
	FnSaveUser         func(u *models.User)
	SaveUserCalled     bool
	SaveUserCalledWith *models.User

	FnGetUserByID         func(id models.UserID) *models.User
	GetUserByIDCalled     bool
	GetUserByIDCalledWith models.UserID

	FnGetUserByEmail         func(email string) *models.User
	GetUserByEmailCalled     bool
	GetUserByEmailCalledWith string
}

func (r *UserRepositoryMock) SaveUser(u *models.User) {
	r.SaveUserCalled = true
	r.SaveUserCalledWith = u

	if r.FnSaveUser != nil {
		r.FnSaveUser(u)
	}
}

func (r *UserRepositoryMock) GetByID(id models.UserID) *models.User {
	r.GetUserByIDCalled = true
	r.GetUserByIDCalledWith = id

	if r.FnGetUserByID != nil {
		return r.FnGetUserByID(id)
	}
	return nil
}

func (r *UserRepositoryMock) GetByEmail(email string) *models.User {
	r.GetUserByEmailCalled = true
	r.GetUserByEmailCalledWith = email

	if r.FnGetUserByEmail != nil {
		return r.FnGetUserByEmail(email)
	}

	return nil
}

type ServerMock struct {
	FnPrepareDisk             func(d *models.Disk) error
	PrepareDiskCalled         bool
	PrepareDiskCalledWithDisk *models.Disk
	PrepareDiskCalledWithUser *models.User

	FnRun     func()
	RunCalled bool
}

func (s *ServerMock) PrepareDisk(d *models.Disk, u *models.User) error {
	s.PrepareDiskCalledWithDisk = d
	s.PrepareDiskCalledWithUser = u
	s.PrepareDiskCalled = true

	if s.FnPrepareDisk != nil {
		return s.FnPrepareDisk(d)
	}

	return nil
}

func (s *ServerMock) Run() {
	s.RunCalled = true

	if s.FnRun != nil {
		s.FnRun()
	}
}
