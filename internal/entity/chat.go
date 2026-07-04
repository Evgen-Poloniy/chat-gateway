package entity

import "time"

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID      int64     `json:"chat_id"`
	SenderID    int64     `json:"sender_id" binding:"required"`
	RecipientID int64     `json:"recipient_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
}
