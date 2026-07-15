package entity

import (
	"time"

	"github.com/google/uuid"
)

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID         uuid.UUID  `validate:"-"`
	ParticipantIDs uuid.UUIDs `validate:"required,len=2,dive,gt=0"`
	CreatedAt      time.Time  `validate:"-"`
}

// GroupChat represents database model for creating group chat.
type GroupChat struct {
	ChatID         uuid.UUID  `validate:"-"`
	ParticipantIDs uuid.UUIDs `validate:"required,min=2,dive,gt=0"`
	Name           string     `validate:"required,min=1,max=64"`
	Title          *string    `validate:"omitempty,max=64"`
	Description    *string    `validate:"omitempty,max=255"`
	CreatedAt      time.Time  `validate:"-"`
	OwnerID        uuid.UUID  `validate:"required"`
}

// Chat represents model of chat
type Chat struct {
	ChatID      uuid.UUID  `validate:"required"`
	ChatType    string     `validate:"required,oneof=direct group,max=16"`
	Name        *string    `validate:"omitempty,max=64"`
	Title       *string    `validate:"omitempty,max=64"`
	Description *string    `validate:"omitempty,max=255"`
	CreatedAt   time.Time  `validate:"-"`
	OwnerID     *uuid.UUID `validate:"omitempty,gt=0"`
}

// UpdateGroupChat represents database model for update data about group chat.
type UpdateGroupChat struct {
	UserIDUpdater uuid.UUID  `validate:"required"`
	ChatID        uuid.UUID  `validate:"required"`
	Name          *string    `validate:"omitempty,min=1,max=64"`
	Title         *string    `validate:"omitempty,max=64"`
	Description   *string    `validate:"omitempty,max=255"`
	CreatedAt     time.Time  `validate:"-"`
	OwnerID       *uuid.UUID `validate:"omitempty,gt=0"`
}

// SendMessage represents model for message broker
type SendMessage struct {
	ChatID   uuid.UUID `validate:"required"`
	SenderID uuid.UUID `validate:"required"`
	Message  string    `validate:"required,min=1,max=4096"`
}
