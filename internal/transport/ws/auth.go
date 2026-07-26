package ws

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMessage represents DTO for websocket connection upgrading
type AuthMessage struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

// extractUUIDFromToken parses, validates access token and extracts user UUID.
func (h *Handler) extractUUIDFromToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(accessToken, h.jwks.Keyfunc, jwt.WithValidMethods(h.authConf.SigningAlgorithms))
	if err != nil {
		return uuid.UUID{}, errs.ErrInvalidToken
	}

	if !token.Valid {
		return uuid.UUID{}, errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, errs.ErrInvalidTokenClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return uuid.UUID{}, errs.ErrInvalidTokenClaims
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.UUID{}, errs.ErrInvalidTokenClaims
	}

	return userID, nil
}
