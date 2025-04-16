package exception

type ValidationError struct {
	Message map[string]string `json:"message"`
}

func NewValidationError(message map[string]string) ValidationError {
	return ValidationError{
		Message: message,
	}
}
