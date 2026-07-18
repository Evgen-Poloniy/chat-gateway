package middleware

import (
	"errors"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/errs"

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
	errs.CodeQueryError:           http.StatusInternalServerError,
	errs.CodeTransactionError:     http.StatusInternalServerError,
	errs.CodeSerializationError:   http.StatusInternalServerError,
	errs.CodeDeserializationError: http.StatusInternalServerError,
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
		if httpError, ok := errors.AsType[*errs.HttpError](err); ok {
			statusCode = httpError.StatusCode
			code = httpError.Code
			message = httpError.Message
		} else if appError, ok := errors.AsType[*errs.AppError](err); ok {
			var ok bool
			if statusCode, ok = httpStatusMap[appError.Code]; !ok {
				statusCode = http.StatusInternalServerError
			}

			if code, ok = errs.MapToString[appError.Code]; !ok {
				code = "unknown_error"
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
