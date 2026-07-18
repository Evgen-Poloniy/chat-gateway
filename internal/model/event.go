package model

import (
	"time"
)

// EventMessage represents message from producer
type EventMessage struct {
	MessageID string    `json:"message_id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// Event represents object from streaming method
type Event struct {
	Data EventMessage
	Err  error
}
