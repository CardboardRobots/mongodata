package repository

type RepositoryError struct {
	Message string
}

func (err RepositoryError) Error() string {
	return err.Message
}

func NewErrNotFound() RepositoryError {
	return RepositoryError{
		Message: "not found",
	}
}
