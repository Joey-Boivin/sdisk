package models

import "github.com/google/uuid"

type UserID struct {
	id uuid.UUID
}

func NewUserID() UserID {
	return UserID{
		id: uuid.New(),
	}
}

func FromString(id string) (UserID, error) {

	userID, err := uuid.Parse(id)

	if err != nil {
		return UserID{}, &ErrInvalidID{invalidID: id}
	}

	return UserID{id: userID}, nil
}

func (u *UserID) ToString() string {
	return u.id.String()
}

type User struct {
	id       UserID
	email    string
	password string
	disk     *Disk
}

func NewUser(email string, password string) *User {
	return &User{
		NewUserID(),
		email,
		password,
		nil,
	}
}

func (u *User) GetID() UserID {
	return u.id
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetDiskSpaceLeft() (uint64, error) {
	if u.disk == nil {
		return 0, &ErrUserHasNoDisk{}
	}

	return u.disk.GetSpaceLeft(), nil
}

func (u *User) AddDisk(d *Disk) error {
	if u.disk != nil {
		return &ErrUserAlreadyHasADisk{}
	}

	u.disk = d
	return nil
}
