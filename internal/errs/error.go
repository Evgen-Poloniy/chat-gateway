package errs

import "errors"

type ErrCode string

const (
	CodeUserNotFound         ErrCode = "user_not_found"
	CodeChatNotFound         ErrCode = "chat_not_found"
	CodeUsersNotFound        ErrCode = "users_not_found"
	CodeUserIsNotInChat      ErrCode = "user_is_not_chat_member"
	CodeUserHaveNotChats     ErrCode = "user_has_no_chats"
	CodeUniqueViolation      ErrCode = "unique_violation"
	CodeForeignKeyViolation  ErrCode = "foreign_key_violation"
	CodeQueryError           ErrCode = "database_query_error"
	CodeTransactionError     ErrCode = "database_transaction_error"
	CodeSerializationError   ErrCode = "serialization_error"
	CodeDeserializationError ErrCode = "deserialization_error"
	CodeUsernameIsRequired   ErrCode = "username_is_required"
	CodeUserIdIsRequired     ErrCode = "user_id_is_required"
	CodeChatIdIsRequired     ErrCode = "chat_id_is_required"
	CodeInvalidParameter     ErrCode = "invalid_parameter"
	CodeValidationError      ErrCode = "validation_error"
	CodeRedisError           ErrCode = "redis_error"
	CodeKafkaError           ErrCode = "kafka_error"
	CodeFailedToExpireKey    ErrCode = "failed_to_expire_key"
	CodeFailedEventChannel   ErrCode = "failed_event_channel"
	CodeEmptyUserIDs         ErrCode = "empty_user_ids"
	CodeDeliveryFailed       ErrCode = "delivery_failed"
	CodeBadRequest           ErrCode = "bad_request"
	CodeMissingAuthHeaders   ErrCode = "missing_auth_headers"
	CodeWrongAuthHeader      ErrCode = "wrong_auth_header"
	CodeInvalidAPIKey        ErrCode = "invalid_api_key"
	CodeInvalidToken         ErrCode = "invalid_token"
	CodeInvalidTokenClaims   ErrCode = "invalid_token_claims"
)

// App error
type AppError struct {
	Code    ErrCode
	Message string
	Err     error
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

// Repository errors.
var (
	ErrFailedEventChannel = errors.New("channel closed unexpectedly")
)

// Validation errors.
var (
	ErrInvalidChatID            = errors.New("invalid chat_id")
	ErrEmptyUserIDs             = errors.New("empty list of user ids")
	ErrUserIDsConversionFailure = errors.New("user ids conversion failure")
	ErrUserIsNotInChat          = errors.New("user is not chat member")
)

// Transport errors.
var (
	ErrMissingAuthHeader      = errors.New("authorization header is required")
	ErrInvalidAuthHeader      = errors.New("invalid authorization header format")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrInvalidTokenClaims     = errors.New("invalid token claims")
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
