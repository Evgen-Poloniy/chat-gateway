package middleware

import (
	"errors"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
)

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
		if httpError, ok := errors.AsType[*errs.HttpError](err); ok {
			statusCode = httpError.StatusCode
			code = httpError.Code
			message = httpError.Message
		} else if appError, ok := errors.AsType[*errs.AppError](err); ok {
			switch appError.Code {
			case errs.CodeUserNotFound:
				statusCode = http.StatusNotFound
				code = "user_not_found"

			case errs.CodeChatNotFound:
				statusCode = http.StatusNotFound
				code = "chat_not_found"

			case errs.CodeUniqueViolation:
				statusCode = http.StatusConflict
				code = "unique_violation"

			case errs.CodeUsernameIsRequired:
				statusCode = http.StatusBadRequest
				code = "username_is_required"

			case errs.CodeUserIdIsRequired:
				statusCode = http.StatusBadRequest
				code = "user_id_is_required"

			case errs.CodeInvalidParameter:
				statusCode = http.StatusBadRequest
				code = "invalid_parameter"

			case errs.CodeQueryError:
				statusCode = http.StatusInternalServerError
				code = "database_query_error"

			case errs.CodeTransactionError:
				statusCode = http.StatusInternalServerError
				code = "database_transaction_error"

			case errs.CodeSerializationError, errs.CodeDeserializationError:
				statusCode = http.StatusInternalServerError
				code = "internal_server_error"

			default:
				statusCode = http.StatusInternalServerError
				code = "unknown_error"
			}

			message = appError.Message
		} else {
			statusCode = http.StatusInternalServerError
			code = "unknown error"
			message = "unknown error"
			c.Error(err)
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
