package middleware

import (
	"strings"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// APIKeyAuth checks request on API-Key availability and validity
func APIKeyAuth(apiKeyHash []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errs.NewAppError(
				errs.CodeMissingAuthHeaders,
				errs.ErrMissingAuthHeader.Error(),
				errs.ErrMissingAuthHeader,
			))
			c.Abort()
			return
		}

		const prefix = "API-KEY "
		if !strings.HasPrefix(authHeader, prefix) {
			c.Error(errs.NewAppError(
				errs.CodeWrongAuthHeader,
				errs.ErrWrongAuthHeader.Error(),
				errs.ErrWrongAuthHeader,
			))
			c.Abort()
			return
		}

		apiKey := strings.TrimPrefix(authHeader, prefix)

		if bcrypt.CompareHashAndPassword(apiKeyHash, []byte(apiKey)) != nil {
			c.Error(errs.NewAppError(
				errs.CodeInvalidAPIKey,
				errs.ErrInvalidAPIKey.Error(),
				errs.ErrInvalidAPIKey,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
