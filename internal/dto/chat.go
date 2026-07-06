package dto

import (
	"time"
)

// CreateGroupChat represents DTO for creating group chat.
type CreateGroupChat struct {
	ParticipantIDs []int64 `json:"participant_ids" binding:"required,min=2"`
	Name           string  `json:"name" binding:"required"`
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	OwnerID        int64   `json:"owner_id" binding:"required"`
}

// CreateDirectChat represents DTO for creating direct chat.
type CreateDirectChat struct {
	SenderID    int64 `json:"sender_id" binding:"required"`
	RecipientID int64 `json:"recipient_id" binding:"required"`
}

// DirectChatData represents DTO for direct chat response data.
type DirectChatData struct {
	ChatID    int64     `json:"chat_id"`
	CreatedAt time.Time `json:"created_at"`
}

// GroupChatData represents DTO for group chat response data.
type GroupChatData struct {
	ChatID         int64     `json:"chat_id"`
	ParticipantIDs []int64   `json:"participant_ids"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Title          *string   `json:"title,omitempty"`
	Description    *string   `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	OwnerID        int64     `json:"owner_id"`
}
