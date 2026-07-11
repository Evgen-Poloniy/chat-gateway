package dto

import (
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// CreateDirectChatReq represents DTO for creating direct chat.
type CreateDirectChatReq struct {
	ParticipantIDs []int64 `json:"participant_ids"`
}

// CreateGroupChatReq represents DTO for creating group chat.
type CreateGroupChatReq struct {
	ParticipantIDs []int64 `json:"participant_ids" binding:"required,min=2"`
	Name           string  `json:"name" binding:"required"`
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	OwnerID        int64   `json:"owner_id" binding:"required"`
}

// DirectChatResp represents DTO for direct chat response data.
type DirectChatResp struct {
	ChatID         int64     `json:"chat_id"`
	ParticipantIDs []int64   `json:"participant_ids"`
	ChatType       string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
}

// GroupChatResp represents DTO for group chat response data.
type GroupChatResp struct {
	ChatID         int64     `json:"chat_id"`
	ParticipantIDs []int64   `json:"participant_ids"`
	ChatType       string    `json:"type"`
	Name           string    `json:"name"`
	Title          *string   `json:"title,omitempty"`
	Description    *string   `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	OwnerID        int64     `json:"owner_id"`
}

// UpdateGroupChatReq represents DTO request for updating data about group chat.
type UpdateGroupChatReq struct {
	UserIDUpdater int64   `json:"user_id_updater"`
	Name          *string `json:"name,omitempty"`
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	OwnerID       *int64  `json:"owner_id,omitempty"`
}

// UpdateGroupChatResp represents DTO response about updating group data.
type UpdateGroupChatResp struct {
	UserIDUpdater int64     `json:"user_id_updater"`
	ChatID        int64     `json:"chat_id"`
	ChatType      string    `json:"type"`
	Name          string    `json:"name"`
	Title         *string   `json:"title,omitempty"`
	Description   *string   `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	OwnerID       int64     `json:"owner_id"`
}

// SendMessage represents DTO for message sending
type SendMessageReq struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	Message  string `json:"message" binding:"required"`
}

// SendMessage represents DTO for message sending
type SendMessageResp struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	Status   string `json:"status"`
}

// ChatsResp represents DTO response with user chats
type ChatsResp struct {
	Chats []entity.Chat `json:"chats"`
}
