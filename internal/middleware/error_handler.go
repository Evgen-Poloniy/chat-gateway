package middleware

import (
	"chat-gateway/internal/dto"
	errs "chat-gateway/pkg/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode := c.Writer.Status()

			if code, exists := c.Get("code"); exists {
				codeStr, ok := code.(string)
				if !ok {
					codeStr = "UNKNOWN_ERROR"
				}

				c.AbortWithStatusJSON(statusCode, dto.ErrorResponse{
					Error: dto.Error{
						Code:    codeStr,
						Message: err.Error(),
					},
				})
			}

			c.AbortWithStatusJSON(statusCode, dto.ErrorResponse{
				Error: dto.Error{
					Code:    "UNKNOWN_ERROR",
					Message: err.Error(),
				},
			})

			return
		}

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var code string
		statusCode := c.Writer.Status()
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
		} else {
			statusCode = http.StatusInternalServerError
			code = "UNKNOWN_ERROR"
			c.Error(errors.New("unknown error"))
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error: dto.Error{
				Code:    code,
				Message: err.Error(),
			},
		})
	}
}
