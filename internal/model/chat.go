package model

import (
	"time"

	"github.com/google/uuid"
)

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID         uuid.UUID
	ParticipantIDs uuid.UUIDs
	CreatedAt      time.Time
}

// GroupChat represents database model for creating direct chat.
type GroupChat struct {
	ChatID         uuid.UUID
	ParticipantIDs uuid.UUIDs
	Name           string
	Title          *string
	Description    *string
	CreatedAt      time.Time
	OwnerID        uuid.UUID
}

// Chat represents model of chat
type Chat struct {
	ChatID      uuid.UUID  `db:"id"`
	ChatType    string     `db:"type"`
	Name        *string    `db:"name"`
	Title       *string    `db:"title"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	OwnerID     *uuid.UUID `db:"owner_id"`
}

// UpdateGroupChat represents database model for update data about group chat.
type UpdateGroupChat struct {
	UserIDUpdater uuid.UUID  `db:"user_id_updater"`
	ChatID        uuid.UUID  `db:"id"`
	Name          *string    `db:"name"`
	Title         *string    `db:"title"`
	Description   *string    `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	OwnerID       *uuid.UUID `db:"owner_id"`
}

// SendMessage represents model for message broker
type SendMessage struct {
	ChatID   uuid.UUID `json:"chat_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Message  string    `json:"message" binding:"required"`
}
