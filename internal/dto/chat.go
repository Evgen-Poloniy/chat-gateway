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

// DispatchMessageReq represents DTO request for dispatch message_id to all online users
type DispatchMessageReq struct {
	MessageID int64 `json:"message_id" binding:"required,min=1"`

	SenderID int64 `json:"sender_id" binding:"required,min=1"`

	// Get online user list temporary. Soon will be replace on chat_id
	RecipientIDs []int64 `json:"recipient_ids" binding:"required, min=1"`

	Message string `json:"message" binding:"required,min=1"`
}

// DispatchMessage represents DTO dispatch message for websocket connection user
type DispatchMessage struct {
	SenderID int64  `json:"sender_id" binding:"required,min=1"`
	Message  string `json:"message" binding:"required,min=1"`
}

// DispatchMessage represents DTO response of dispatch message to user
type DispatchMessageResp struct {
	OfflineUserIDs []int64 `json:"offline_user_ids"`
}
