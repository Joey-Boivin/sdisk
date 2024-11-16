package models

type UserRepository interface {
	SaveUser(u *User)
	GetUser(id string) *User
}

type User struct {
	email    string
	password string
}

func NewUser(email string, password string) *User {
	return &User{
		email,
		password,
	}
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}
