package exception

type NotFoundError struct {
	Message string `json:"message"`
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Message: message,
	}
}
