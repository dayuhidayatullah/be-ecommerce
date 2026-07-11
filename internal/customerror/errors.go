package customerror

import "fmt"

// ErrorType defines the category of the error
type ErrorType string

const (
	TypeNotFound     ErrorType = "NOT_FOUND"
	TypeConflict     ErrorType = "CONFLICT"
	TypeBadRequest   ErrorType = "BAD_REQUEST"
	TypeUnauthorized ErrorType = "UNAUTHORIZED"
	TypeInternal     ErrorType = "INTERNAL_SERVER_ERROR"
)

// AppError represents a custom error that holds the error type and message
type AppError struct {
	Type    ErrorType
	Message string
	Err     error // original underlying error if any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Helper functions to create specific errors easily
func NewNotFoundError(message string) *AppError {
	return &AppError{Type: TypeNotFound, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Type: TypeConflict, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Type: TypeBadRequest, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Type: TypeUnauthorized, Message: message}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{Type: TypeInternal, Message: message, Err: err}
}
