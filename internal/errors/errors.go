package errors

import "fmt"

type ErrorType string

const (
	NotFound      ErrorType = "NOT_FOUND"
	BadRequest    ErrorType = "BAD_REQUEST"
	Internal      ErrorType = "INTERNAL"
	Validation    ErrorType = "VALIDATION"
	AlreadyExists ErrorType = "ALREADY_EXISTS"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAlreadyExists(message string, err error) *Error {
	return &Error{
		Type:    AlreadyExists,
		Message: message,
		Err:     err,
	}
}
func NewInternal(message string, err error) *Error {
	return &Error{
		Type:    Internal,
		Message: message,
		Err:     err,
	}
}
