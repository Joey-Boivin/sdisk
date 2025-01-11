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

type ErrUserHasNoDisk struct {
}

func (e *ErrUserHasNoDisk) Error() string {
	return "user has no disk"
}

type ErrInvalidID struct {
	invalidID string
}

func (e *ErrInvalidID) Error() string {
	return fmt.Sprintf("the id %s is not valid", e.invalidID)
}
