package ws

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents websocket connection with user.
type Client struct {
	userID uuid.UUID
	conn   *websocket.Conn
	send   chan dto.DispatchMessage
}

func NewClient(userID uuid.UUID, conn *websocket.Conn) *Client {
	return &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan dto.DispatchMessage, 64),
	}
}
