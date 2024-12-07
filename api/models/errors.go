package models

import (
	"fmt"
)

type ErrUserAlreadyHasADisk struct {
	Email string
}

func (e *ErrUserAlreadyHasADisk) Error() string {
	return fmt.Sprintf("user already exists for the email address %s", e.Email)
}
