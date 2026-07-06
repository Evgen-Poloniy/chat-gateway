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
			err := c.Errors.Last().Err
			statusCode := c.MustGet("status_code").(int)
			code := c.MustGet("code").(string)

			c.AbortWithStatusJSON(statusCode, dto.ResponseError{
				Error: dto.Error{
					Code:    code,
					Message: err.Error(),
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
		err := c.Errors.Last().Err

		// Error mapping
		if appError, ok := errors.AsType[*errs.AppError](err); ok {
			statusCode = appError.StatusCode
			code = appError.Code
		} else if c.Writer.Status() == http.StatusNotFound {
			statusCode = http.StatusNotFound
			code = "NOT_FOUND"
		} else if c.Writer.Status() == http.StatusMethodNotAllowed {
			statusCode = http.StatusMethodNotAllowed
			code = "METHOD_NOT_ALLOWED"
		} else if errors.Is(err, errs.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			code = "RECORD_NOT_FOUND"
		} else {
			statusCode = http.StatusInternalServerError
			code = "UNKNOWN_ERROR"
			c.Error(errors.New("unknown error"))
		}

		c.Set("status_code", statusCode)
		c.Set("code", code)

		c.JSON(statusCode, dto.ResponseError{
			Error: dto.Error{
				Code:    code,
				Message: err.Error(),
			},
		})
	}
}
