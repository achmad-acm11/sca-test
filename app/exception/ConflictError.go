package exception

type ConflictError struct {
	Message string `json:"message"`
}

func NewConflictError(message string) ConflictError {
	return ConflictError{
		Message: message,
	}
}
