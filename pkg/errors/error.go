package errs

import "errors"

type ErrCode int

const (
	CodeUserNotFound ErrCode = iota
	CodeChatNotFound
	CodeUserHaveNotChats
	CodeUniqueViolation
	CodeQueryError
	CodeTransactionError
	CodeSerializationError
	CodeDeserializationError
	CodeUsernameIsRequired
	CodeUserIdIsRequired
	CodeInvalidParameter
)

// App error
type AppError struct {
	Code    ErrCode
	Message string
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
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(code ErrCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
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

// Transport errors.
var (
	ErrUsernameIsRequired     = errors.New("username is required")
	ErrUserIdIsRequired       = errors.New("user_id is required")
	ErrInvalidParameter       = errors.New("invalid parameter")
	ErrPageRequiredBeGreater  = errors.New("query parameter 'page' required be greater or equals then 1")
	ErrLimitRequiredBeGreater = errors.New("query parameter 'limit' required be greater or equals then 1")
	ErrLimitRequiredBeLess    = errors.New("query parameter 'limit' required be less or equals then 100")
)
