package mocks

type RegisterService struct {
	RegisterUserIndirection func(email string, password string) error
	RegisterUserInvoked     bool
}

func (r *RegisterService) RegisterUser(email string, password string) error {
	r.RegisterUserInvoked = true
	if r.RegisterUserIndirection != nil {
		return r.RegisterUserIndirection(email, password)
	}
	return nil
}
