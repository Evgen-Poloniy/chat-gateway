package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	ErrSetDeadline = errors.New("error when set websocket read deadline")
)

// HandleWebSocketUpgrade upgrades an HTTP connection to a WebSocket connection.
func (h *Handler) WebSocketUpgrade(c *gin.Context) {
	requestID := uuid.NewString()
	start := time.Now()

	entry := h.logger.WithFields(map[string]interface{}{
		"id":     requestID,
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})
	entry.Info("websocket upgrade request start")

	conn, err := h.wsHub.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		message := fmt.Sprintf("failed to upgrade connection to websocket: %v", err)
		abortWithError(c, h.logger, http.StatusInternalServerError, "failed_to_upgrade_to_websocket_connection", message, start, requestID)
		return
	}

	userID, err := h.getUserID(conn, time.Second*5)
	if err != nil {
		defer func() {
			if err := conn.Close(); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					logError(h.logger, h.wsHub.serverAddress, err)
				}
			}
		}()

		message := fmt.Sprintf("authentication error: %v", err)

		if errors.Is(err, ErrSetDeadline) {
			abortWithError(c, h.logger, http.StatusInternalServerError, "internal_server_error", message, start, requestID)
			return
		}

		abortWithError(c, h.logger, http.StatusUnauthorized, "unauthorized", message, start, requestID)
		return
	}

	client := NewClient(userID, conn, h.wsHub.dispatchConf.SendChannelCap)

	h.wsHub.register(client)

	defer func() {
		if err := h.wsHub.unregister(client); err != nil {
			h.logger.WithFields(map[string]interface{}{
				"id":         requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"close_code": websocket.CloseInternalServerErr,
				"code":       "failed_to_close_websocket_connection",
			}).Error(fmt.Errorf("failed to close websocket connection: %w", err))
		}
	}()

	go h.wsHub.writeMessage(client)

	h.wsHub.readMessage(client)
}

// getUserID waits access token from client by websocket connection and extract user_id.
func (h *Handler) getUserID(conn *websocket.Conn, timeout time.Duration) (uuid.UUID, error) {
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return uuid.UUID{}, ErrSetDeadline
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return uuid.UUID{}, errors.New("timeout or read error")
	}

	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return uuid.UUID{}, ErrSetDeadline
	}

	var msg AuthMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return uuid.UUID{}, errors.New("invalid json format")
	}

	if msg.Type != "auth" || msg.AccessToken == "" {
		return uuid.UUID{}, errors.New("missing access token")
	}

	userID, err := h.extractUUIDFromToken(msg.AccessToken)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

// abortWithError logs the error and sends a formatted JSON error response to the client
func abortWithError(c *gin.Context, logger *logrus.Logger, statusCode int, code, message string, start time.Time, requestID string) {
	logger.WithFields(map[string]interface{}{
		"id":          requestID,
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"ip":          c.ClientIP(),
		"latency":     time.Since(start),
		"status_code": statusCode,
		"code":        code,
	}).Error(message)

	c.AbortWithStatusJSON(
		statusCode,
		dto.ResponseError{
			Error: dto.Error{
				Code:    code,
				Message: message,
			},
		},
	)
}
