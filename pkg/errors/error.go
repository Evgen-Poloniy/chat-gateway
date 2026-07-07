package errs

import "errors"

// App error
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Http error type
type HttpError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(statusCode int, code, message string, err error) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// Repository errors.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Service errors.
var (
	ErrChatAlreadyExists = errors.New("chat already exists")
)

// Transport errors.
var (
	ErrUsernameIsRequired = errors.New("username is required")
	ErrUserIdIsRequired   = errors.New("user_id is required")
	ErrInvalidParameter   = errors.New("invalid parameter")
)
