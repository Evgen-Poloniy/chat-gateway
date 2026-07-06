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
			statusCode = http.StatusInternalServerError
			code = appError.Code
			message = appError.Message
			c.Error(appError.Err)
		} else if errors.Is(err, errs.ErrRecordNotFound) {
			statusCode = http.StatusNotFound
			code = "RECORD_NOT_FOUND"
			message = err.Error()
		} else {
			statusCode = http.StatusInternalServerError
			code = "UNKNOWN_ERROR"
			message = "unknown error"
			c.Error(errors.New("unknown error"))
		}

		c.Set("status_code", statusCode)
		c.Set("code", code)

		c.JSON(statusCode, dto.ResponseError{
			Error: dto.Error{
				Code:    code,
				Message: message,
			},
		})
	}
}
