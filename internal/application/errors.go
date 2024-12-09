package application

import (
	"fmt"
)

type ErrUserAlreadyExists struct {
	Email string
}

func (e *ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("user already exists for the email address %s", e.Email)
}

type ErrUserDoesNotExist struct {
	Email string
}

func (e *ErrUserDoesNotExist) Error() string {
	return fmt.Sprintf("user with email %s does not exist", e.Email)
}
