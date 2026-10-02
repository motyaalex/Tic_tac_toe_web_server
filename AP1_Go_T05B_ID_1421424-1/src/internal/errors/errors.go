package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrValidation   = errors.New("validation error")
)

type AppError struct {
	Err     error
	Message string
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewNotFound(msg string) error {
	return &AppError{Err: ErrNotFound, Message: msg}
}
