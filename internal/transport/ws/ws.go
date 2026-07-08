package ws

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Hub manages WebSocket connections
type Hub struct {
	upgrader websocket.Upgrader
	mu       sync.RWMutex
	conns    map[int]*websocket.Conn
}

var hub = &Hub{
	upgrader: websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	},
	conns: make(map[int]*websocket.Conn),
}

// HandleWebSocketUpgrade upgrades an HTTP connection to a WebSocket connection.
func WebSocketUpgrade(c *gin.Context, logger *logrus.Logger) {
	id := uuid.New().ID()
	start := time.Now()

	entry := logger.WithFields(map[string]interface{}{
		"id":     id,
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})
	entry.Info("websocket upgrade request start")

	queryUserID := c.Query("user_id")
	if queryUserID == "" {
		abortWithError(c, logger, http.StatusBadRequest, "bad_request", "user_id is required", start, id)
		return
	}

	userID, err := strconv.Atoi(queryUserID)
	if err != nil {
		abortWithError(c, logger, http.StatusBadRequest, "bad_request", "user_id must be an integer", start, id)
		return
	}

	conn, err := hub.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		message := fmt.Sprintf("failed to upgrade connection to websocket: %v", err)
		abortWithError(c, logger, http.StatusInternalServerError, "failed_to_upgrade_to_websocket_connection", message, start, id)
		return
	}

	hub.mu.Lock()
	hub.conns[userID] = conn
	hub.mu.Unlock()

	defer func() {
		hub.mu.Lock()
		delete(hub.conns, userID)
		hub.mu.Unlock()

		if err := conn.Close(); err != nil {
			logger.WithFields(map[string]interface{}{
				"id":         id,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"close_code": websocket.CloseInternalServerErr,
				"code":       "failed_to_close_websocket_connection",
			}).Error(fmt.Errorf("failed to close websocket connection: %w", err))
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			var closeCode int
			var code string
			var message string
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				closeCode = closeErr.Code
			}

			switch closeCode {
			case websocket.CloseNormalClosure:
				closeCode = websocket.CloseNormalClosure
				message = "websocket connection closed normally"
			case websocket.CloseGoingAway:
				closeCode = websocket.CloseGoingAway
				message = "websocket connection going away"
			case websocket.CloseAbnormalClosure:
				closeCode = websocket.CloseAbnormalClosure
				code = "websocket_connection_unexpectedly_closed"
				message = fmt.Sprintf("websocket connection closed unexpectedly: %v", err)
			case websocket.CloseInternalServerErr:
				closeCode = websocket.CloseInternalServerErr
				code = "websocket_internal_server_error"
				message = fmt.Sprintf("websocket connection error: %v", err)
			default:
				closeCode = websocket.CloseInternalServerErr
				code = "websocket_unknown_error"
				message = fmt.Sprintf("websocket connection error: %v", err)
			}

			switch closeCode {
			case websocket.CloseNormalClosure, websocket.CloseGoingAway:
				logger.WithFields(map[string]interface{}{
					"id":         id,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"ip":         c.ClientIP(),
					"close_code": closeCode,
				}).Info(message)
			default:
				logger.WithFields(map[string]interface{}{
					"id":         id,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"ip":         c.ClientIP(),
					"close_code": closeCode,
					"code":       code,
				}).Error(message)
			}

			break
		}
	}
}

// abortWithError logs the error and sends a formatted JSON error response to the client
func abortWithError(c *gin.Context, logger *logrus.Logger, statusCode int, code, message string, start time.Time, id uint32) {
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
