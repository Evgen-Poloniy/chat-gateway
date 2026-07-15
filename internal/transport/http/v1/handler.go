package v1

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/transport/http/ws"
	"github.com/google/uuid"
)

// MessengerService represents interface for work with messenger business logic.
type MessengerService interface {
	// Users:

	// CreateUser allows create user into messenger by template.
	CreateUser(ctx context.Context, user *entity.User) error

	// GetUserDataByUsername represents searching all data about user from the messenger database.
	GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetUserDataByUserID gets all data about user from the messenger database by user_id.
	GetUserDataByUserID(ctx context.Context, userID uuid.UUID) (*entity.User, error)

	// GetUserIDsByChatID gets user_id by all users who are in the chat.
	GetUserIDsByChatID(ctx context.Context, chatID uuid.UUID) (uuid.UUIDs, error)

	// UpdateUser updates data about user into messenger database.
	UpdateUser(ctx context.Context, user *entity.UpdateUser) error

	// Chats:

	// CreateDirectChat accept user IDs and create direct chat between two users.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error

	// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
	CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error

	// GetChatsByUserID gets chat by user_id with limits and offset.
	GetChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Chat, error)

	// UpdateGroupChat updates data about chat like name, title, description, owner.
	UpdateGroupChat(ctx context.Context, chat *entity.UpdateGroupChat) error

	// SendMessage sends message into target chat.
	SendMessage(ctx context.Context, message *entity.SendMessage) error
}

type Handler struct {
	messenger MessengerService
	wsHub     *ws.Hub
}

func NewHandler(messenger MessengerService, wsHub *ws.Hub) *Handler {
	return &Handler{
		messenger: messenger,
		wsHub:     wsHub,
	}
}
