package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// APIKeyAuth checks request on API-Key availability and validity.
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.Error(errs.NewAppError(
				errs.CodeMissingAuthHeaders,
				errs.ErrMissingAuthHeader.Error(),
				errs.ErrMissingAuthHeader,
			))
			c.Abort()
			return
		}

		if key != apiKey {
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

// InitKeyfunc starts background worker for update IdP public key fro parse token.
func InitKeyfunc(jwksURL string, logger *logrus.Logger) (keyfunc.Keyfunc, error) {
	override := keyfunc.Override{
		RefreshInterval: time.Hour,

		RefreshErrorHandlerFunc: func(u string) func(ctx context.Context, err error) {
			return func(ctx context.Context, err error) {
				logger.Errorf("JWKS refresh error for URL %s: %v", u, err)
			}
		},
		RefreshUnknownKID: rate.NewLimiter(rate.Every(5*time.Minute), 1),
		RateLimitWaitMax:  time.Minute,
	}

	jwks, err := keyfunc.NewDefaultOverrideCtx(context.Background(), []string{jwksURL}, override)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyfunc: %w", err)
	}

	return jwks, nil
}

// AuthMiddleware parses, validates access token and gets token claims.
func AuthMiddleware(jwks keyfunc.Keyfunc, config *config.AuthConfig) gin.HandlerFunc {
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

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.Error(errs.NewAppError(
				errs.CodeWrongAuthHeader,
				errs.ErrInvalidAuthHeader.Error(),
				errs.ErrInvalidAuthHeader,
			))
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, prefix)

		token, err := jwt.Parse(accessToken, jwks.Keyfunc, jwt.WithValidMethods(config.SigningAlgorithms))
		if err != nil {
			c.Error(errs.NewAppError(
				errs.CodeInvalidToken,
				errs.ErrInvalidToken.Error(),
				fmt.Errorf("failed to parse access token: %w", err),
			))
			c.Abort()
			return
		}

		if !token.Valid {
			c.Error(errs.NewAppError(
				errs.CodeInvalidToken,
				errs.ErrInvalidToken.Error(),
				errs.ErrInvalidToken,
			))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(errs.NewAppError(
				errs.CodeInvalidTokenClaims,
				errs.ErrInvalidTokenClaims.Error(),
				errs.ErrInvalidTokenClaims,
			))
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			c.Error(errs.NewAppError(
				errs.CodeInvalidTokenClaims,
				errs.ErrInvalidTokenClaims.Error(),
				fmt.Errorf("invalid token claims: sub"),
			))
			c.Abort()
			return
		}

		var roles []string

		if rawRoles, ok := claims["roles"].([]interface{}); ok {
			for _, r := range rawRoles {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}

		c.Set("user_id", sub)
		c.Set("roles", roles)

		c.Next()
	}
}
