package errs

import (
	"errors"
)

type ErrCode int

const (
	CodeUserNotFound ErrCode = iota
	CodeChatNotFound
	CodeUsersNotFound
	CodeUserHaveNotChats
	CodeUniqueViolation
	CodeForeignKeyViolation
	CodeQueryError
	CodeTransactionError
	CodeSerializationError
	CodeDeserializationError
	CodeUsernameIsRequired
	CodeUserIdIsRequired
	CodeChatIdIsRequired
	CodeInvalidParameter
	CodeValidationError
	CodeRedisError
	CodeFailedToExpireKey
	CodeFailedEventChannel
	CodeEmptyUserIDs
	ErrCodeDeliveryFailed
)

var MapToString = map[ErrCode]string{
	CodeUserNotFound:         "user_not_found",
	CodeChatNotFound:         "chat_not_found",
	CodeUsersNotFound:        "users_not_found",
	CodeUserHaveNotChats:     "user_has_no_chats",
	CodeUniqueViolation:      "unique_violation",
	CodeForeignKeyViolation:  "foreign_key_violation",
	CodeQueryError:           "database_query_error",
	CodeTransactionError:     "database_transaction_error",
	CodeSerializationError:   "internal_server_error",
	CodeDeserializationError: "internal_server_error",
	CodeUsernameIsRequired:   "username_is_required",
	CodeUserIdIsRequired:     "user_id_is_required",
	CodeChatIdIsRequired:     "chat_id_is_required",
	CodeInvalidParameter:     "invalid_parameter",
	CodeValidationError:      "validation_error",
	CodeRedisError:           "redis_error",
	CodeFailedToExpireKey:    "failed_to_expire_key",
	CodeFailedEventChannel:   "failed_event_channel",
	CodeEmptyUserIDs:         "empty_user_ids",
	ErrCodeDeliveryFailed:    "delivery_failed",
}

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

// Repository errors.
var (
	ErrFailedEventChannel = errors.New("redis pubsub channel closed unexpectedly")
)

// Validation errors.
var (
	ErrInvalidChatID            = errors.New("invalid chat_id")
	ErrEmptyUserIDs             = errors.New("empty list of user ids")
	ErrUserIDsConversionFailure = errors.New("user ids conversion failure")
)

// Transport errors.
var (
	ErrMissingAuthHeader      = errors.New("missing authentication header")
	ErrWrongAuthHeader        = errors.New("wrong authentication header")
	ErrInvalidAPIKey          = errors.New("invalid API-Key")
	ErrUsernameIsRequired     = errors.New("username is required")
	ErrUserIdIsRequired       = errors.New("user_id is required")
	ErrChatIdIsRequired       = errors.New("chat_id is required")
	ErrInvalidParameter       = errors.New("invalid parameter")
	ErrPageRequiredBeGreater  = errors.New("query parameter 'page' required be greater or equals then 1")
	ErrLimitRequiredBeGreater = errors.New("query parameter 'limit' required be greater or equals then 1")
	ErrLimitRequiredBeLess    = errors.New("query parameter 'limit' required be less or equals then 100")
	ErrUserIsOffline          = errors.New("user is offline")
)
