package ws

import (
	"net/http"
	"sync"

	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/gorilla/websocket"
)

// wsHub manages WebSocket connections
type Hub struct {
	upgrader websocket.Upgrader
	mu       sync.RWMutex
	conns    map[int64]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		conns: make(map[int64]*websocket.Conn),
	}
}

// register is used only in WebSocketUpgrade for add new ws connection with user
func (h *Hub) register(userID int64, conn *websocket.Conn) {
	h.mu.Lock()
	h.conns[userID] = conn
	h.mu.Unlock()
}

// unregister is used only in WebSocketUpgrade for close ws connection user with user
func (h *Hub) unregister(userID int64, conn *websocket.Conn) error {
	h.mu.Lock()
	delete(h.conns, userID)
	h.mu.Unlock()

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func (h *Hub) SendToUser(userID int64, data interface{}) error {
	h.mu.RLock()
	conn, online := h.conns[userID]
	h.mu.RUnlock()

	if !online {
		return errs.ErrUserIsOffline
	}

	return conn.WriteJSON(data)
}
