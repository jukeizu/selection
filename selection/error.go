package selection

import "fmt"

func NewValidationError(format string, a ...interface{}) ValidationError {
	return ValidationError{Message: fmt.Sprintf(format, a...)}
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
