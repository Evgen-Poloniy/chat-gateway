package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// HandleWebSocketUpgrade upgrades an HTTP connection to a WebSocket connection.
func (h *Handler) WebSocketUpgrade(c *gin.Context) {
	id := uuid.NewString()
	start := time.Now()

	entry := h.logger.WithFields(map[string]interface{}{
		"id":     id,
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})
	entry.Info("websocket upgrade request start")

	queryUserID := c.Query("user_id")
	if queryUserID == "" {
		abortWithError(c, h.logger, http.StatusBadRequest, "bad_request", "user_id is required", start, id)
		return
	}

	userID, err := uuid.Parse(queryUserID)
	if err != nil {
		abortWithError(c, h.logger, http.StatusBadRequest, "bad_request", "user_id must be a valid UUID", start, id)
		return
	}

	conn, err := h.wsHub.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		message := fmt.Sprintf("failed to upgrade connection to websocket: %v", err)
		abortWithError(c, h.logger, http.StatusInternalServerError, "failed_to_upgrade_to_websocket_connection", message, start, id)
		return
	}

	client := NewClient(userID, conn)

	h.wsHub.register(client)

	defer func() {
		if err := h.wsHub.unregister(client); err != nil {
			h.logger.WithFields(map[string]interface{}{
				"id":         id,
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

// abortWithError logs the error and sends a formatted JSON error response to the client
func abortWithError(c *gin.Context, logger *logrus.Logger, statusCode int, code, message string, start time.Time, id string) {
	logger.WithFields(map[string]interface{}{
		"id":          id,
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"ip":          c.ClientIP(),
		"latency":     time.Since(start),
		"status_code": statusCode,
		"code":        code,
	}).Warn(message)

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
