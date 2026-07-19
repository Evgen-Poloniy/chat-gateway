package middleware

import (
	"errors"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"

	"github.com/gin-gonic/gin"
)

var httpStatusMap = map[errs.ErrCode]int{
	errs.CodeUserNotFound:         http.StatusNotFound,
	errs.CodeChatNotFound:         http.StatusNotFound,
	errs.CodeUsersNotFound:        http.StatusNotFound,
	errs.CodeUserHaveNotChats:     http.StatusNotFound,
	errs.CodeUniqueViolation:      http.StatusConflict,
	errs.CodeForeignKeyViolation:  http.StatusConflict,
	errs.CodeUsernameIsRequired:   http.StatusBadRequest,
	errs.CodeUserIdIsRequired:     http.StatusBadRequest,
	errs.CodeChatIdIsRequired:     http.StatusBadRequest,
	errs.CodeInvalidParameter:     http.StatusBadRequest,
	errs.CodeValidationError:      http.StatusBadRequest,
	errs.CodeEmptyUserIDs:         http.StatusBadRequest,
	errs.CodeBadRequest:           http.StatusBadRequest,
	errs.CodeQueryError:           http.StatusInternalServerError,
	errs.CodeTransactionError:     http.StatusInternalServerError,
	errs.CodeSerializationError:   http.StatusInternalServerError,
	errs.CodeDeserializationError: http.StatusInternalServerError,
	errs.CodeRedisError:           http.StatusInternalServerError,
	errs.CodeFailedToExpireKey:    http.StatusInternalServerError,
	errs.CodeFailedEventChannel:   http.StatusInternalServerError,
	errs.CodeDeliveryFailed:       http.StatusInternalServerError,
	errs.CodeMissingAuthHeaders:   http.StatusUnauthorized,
	errs.CodeWrongAuthHeader:      http.StatusUnauthorized,
	errs.CodeInvalidAPIKey:        http.StatusUnauthorized,
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Errors) > 0 {
			statusCode := c.MustGet("status_code").(int)
			code := c.MustGet("code").(string)
			message := c.MustGet("message").(string)

			c.AbortWithStatusJSON(statusCode, dto.ResponseError{
				Error: dto.Error{
					Code:    code,
					Message: message,
				},
			})

			return
		}

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var statusCode int
		var code string
		var message string
		err := c.Errors.Last().Err

		// Error mapping
		if appError, ok := errors.AsType[*errs.AppError](err); ok {
			code = string(appError.Code)

			var ok bool
			if statusCode, ok = httpStatusMap[appError.Code]; !ok {
				statusCode = http.StatusInternalServerError
			}

			message = appError.Message
		} else {
			statusCode = http.StatusInternalServerError
			code = "unknown_error"
			message = "unknown error"
		}

		c.Status(statusCode)
		c.Set("code", code)

		c.JSON(statusCode, dto.ResponseError{
			Error: dto.Error{
				Code:    code,
				Message: message,
			},
		})
	}
}
