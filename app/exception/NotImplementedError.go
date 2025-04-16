package exception

type NotImplementedError struct {
	Message string `json:"message"`
}

func NewNotImplementedError(message string) NotImplementedError {
	return NotImplementedError{
		Message: message,
	}
}
