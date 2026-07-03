package errs

import "errors"

// Error type
type AppError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

// Repository errors
var (
	ErrRecordNotFound = errors.New("record not found")
)
