package v1

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// MessengerService represents interface for work with messenger business logic.
type MessengerService interface {
	// CreateUser allows create user into messenger by template.
	CreateUser(ctx context.Context, user *entity.User) error

	// GetUserDataByUsername represents searching all data about user from the messenger database.
	GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error)

	//GetUserDataByUserID gets all data about user from the messenger database by user_id.
	GetUserDataByUserID(ctx context.Context, userID int64) (*entity.User, error)

	// UpdateUser updates data about user into messenger database.
	UpdateUser(ctx context.Context, user *entity.UpdateUser) error

	// CreateDirectChat accept user IDs and create direct chat between two users.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error

	// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
	CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error

	// SendMessage sends message into target chat.
	SendMessage(message *dto.SendMessageReq) error
}

type Handler struct {
	messenger MessengerService
}

func NewHandler(messenger MessengerService) *Handler {
	return &Handler{
		messenger: messenger,
	}
}
