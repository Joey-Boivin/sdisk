package vcs

type ErrVcsAlreadyInitialized struct {
}

func (e *ErrVcsAlreadyInitialized) Error() string {
	return "already a vcs directory"
}
