package middleware

import (
	"net/http"
	"strings"

	"github.com/Evgen-Poloniy/chat-gateway/pkg/errs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// APIKeyAuth checks request on API-Key availability and validity
func APIKeyAuth(apiKeyHash []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errs.NewHttpError(
				"MISSING_AUTH_HEADERS",
				http.StatusUnauthorized,
				errs.ErrMissingAuthHeader.Error(),
				errs.ErrMissingAuthHeader,
			))
			c.Abort()
			return
		}

		const prefix = "API-KEY "
		if !strings.HasPrefix(authHeader, prefix) {
			c.Error(errs.NewHttpError(
				"WRONG_AUTH_HEADER",
				http.StatusUnauthorized,
				errs.ErrWrongAuthHeader.Error(),
				errs.ErrWrongAuthHeader,
			))
			c.Abort()
			return
		}

		apiKey := strings.TrimPrefix(authHeader, prefix)

		if bcrypt.CompareHashAndPassword(apiKeyHash, []byte(apiKey)) != nil {
			c.Error(errs.NewHttpError(
				"INVALID_API_KEY",
				http.StatusUnauthorized,
				errs.ErrInvalidAPIKey.Error(),
				errs.ErrInvalidAPIKey,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
