package errs

import "errors"

// App error
type AppError struct {
	Code    string
	Message string
	Details string
	Err     error
}

// Http error type
type HttpError struct {
	Code       string
	StatusCode int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message, details string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(code string, statusCode int, message string, err error) *HttpError {
	return &HttpError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// Repository errors.
var (
	ErrQuery                     = errors.New("database query error")
	ErrRecordNotFound            = errors.New("record not found")
	ErrUniqueViolation           = errors.New("inserted value must be unique")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction")
	ErrSerialization             = errors.New("serialization error")
	ErrDeserialization           = errors.New("deserialization error")
)

// Transport errors.
var (
	ErrUsernameIsRequired = errors.New("username is required")
	ErrUserIdIsRequired   = errors.New("user_id is required")
	ErrInvalidParameter   = errors.New("invalid parameter")
)
