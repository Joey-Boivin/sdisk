package models

type User struct {
	email    string
	password string
	disk     *Disk
}

func NewUser(email string, password string) *User {
	return &User{
		email,
		password,
		nil,
	}
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
