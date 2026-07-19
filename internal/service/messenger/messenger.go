package messenger

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/messenger_mocks.go -package=mock_messenger_repository

// MessengerRepository represents interface for work with the messenger database.
type MessengerRepository interface {
	// Users:

	// CreateUser allows create user into messenger database.
	CreateUser(ctx context.Context, user *model.User) error

	// GetUserDataByUsername represents searching all data about user from the messenger database.
	GetUserDataByUsername(ctx context.Context, username string) (*model.User, error)

	//GetUserDataByUserID gets all data about user from the messenger database by user_id.
	GetUserDataByUserID(ctx context.Context, userID uuid.UUID) (*model.User, error)

	// GetUserIDsByChatID gets user_id by all users who are in the chat.
	GetUserIDsByChatID(ctx context.Context, chatID uuid.UUID) (uuid.UUIDs, error)

	// UpdateUser updates data about user into messenger database.
	UpdateUser(ctx context.Context, user *model.UpdateUser) error

	// Chats:

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *model.DirectChat) error

	// CreateChat accept user IDs and create group chat.
	CreateGroupChat(ctx context.Context, chat *model.GroupChat) error

	// GetChatsByUserID gets chat by user_id with limits and offset.
	GetChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Chat, error)

	// UpdateGroupChat updates data about chat like name, title, description, owner.
	UpdateGroupChat(ctx context.Context, chat *model.UpdateGroupChat) error
}

// MessageBroker represents interface for work with the messenger broker.
type MessageBroker interface {
	// SendMessage sends message into target chat.
	SendMessage(ctx context.Context, message *model.SendMessage) error
}

// MessengerCache represents interface for work with the messenger cache.
type MessengerCache interface {
	// AddChatMembers writes user_ids at cache by key chat_id.
	AddChatMembers(ctx context.Context, chatID string, userIDs []string) error

	// RemoveChatUser removes user from Redis by chat_id.
	RemoveChatUser(ctx context.Context, chatID string, userID string) error

	// ExpireChatID expire key TTL
	ExpireChatID(ctx context.Context, chatID string) error

	// GetChatMembers allows get user_ids from cache by chat_id.
	GetChatMembers(ctx context.Context, chatID string) ([]string, error)
}

// MessengerService represents implementation of messenger interface.
type MessengerService struct {
	validate            *validator.Validate
	messengerRepository MessengerRepository
	messageBroker       MessageBroker
	messengerCache      MessengerCache
}

func NewMessengerService(
	messengerRepository MessengerRepository,
	messageBroker MessageBroker,
	messengerCache MessengerCache,
) *MessengerService {
	return &MessengerService{
		validate:            validator.New(),
		messengerRepository: messengerRepository,
		messageBroker:       messageBroker,
		messengerCache:      messengerCache,
	}
}

func (m *MessengerService) validateData(ctx context.Context, data any) error {
	if err := m.validate.StructCtx(ctx, data); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			messages := make([]string, 0, len(validationErrors))
			errMessages := make([]string, 0, len(validationErrors))

			for _, fieldErr := range validationErrors {
				var baseMsg string
				if fieldErr.Param() != "" {
					baseMsg = fmt.Sprintf("field=%s | tag=%s | value='%v' | expected=%s", fieldErr.Field(), fieldErr.Tag(), fieldErr.Value(), fieldErr.Param())
				} else {
					baseMsg = fmt.Sprintf("field=%s | tag=%s | value='%v'", fieldErr.Field(), fieldErr.Tag(), fieldErr.Value())
				}

				messages = append(messages, fmt.Sprintf("[%s]", baseMsg))
				errMessages = append(errMessages, fmt.Sprintf("[%s | err=%s]", baseMsg, fieldErr.Error()))
			}

			return errs.NewAppError(
				errs.CodeValidationError,
				fmt.Sprintf("validation failed: %s", strings.Join(messages, " | ")),
				fmt.Errorf("validation error details: %s", strings.Join(errMessages, " | ")),
			)
		}

		return errs.NewAppError(
			errs.CodeValidationError,
			"validation error",
			fmt.Errorf("validation error: %v", err),
		)
	}

	return nil
}
