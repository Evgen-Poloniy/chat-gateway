package entity

import "time"

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID         int64     `validate:"required,gt=0"`
	ParticipantIDs []int64   `validate:"required,len=2,dive,gt=0"`
	CreatedAt      time.Time `validate:"structonly"`
}

// GroupChat represents database model for creating group chat.
type GroupChat struct {
	ChatID         int64     `validate:"required,gt=0"`
	ParticipantIDs []int64   `validate:"required,min=2,dive,gt=0"`
	Name           string    `validate:"required,min=1,max=64"`
	Title          *string   `validate:"omitempty,max=64"`
	Description    *string   `validate:"omitempty,max=255"`
	CreatedAt      time.Time `validate:"structonly"`
	OwnerID        int64     `validate:"required,gt=0"`
}

// Chat represents model of chat
type Chat struct {
	ChatID      int64     `validate:"required,gt=0"`
	ChatType    string    `validate:"required,oneof=direct group,max=16"`
	Name        *string   `validate:"omitempty,max=64"`
	Title       *string   `validate:"omitempty,max=64"`
	Description *string   `validate:"omitempty,max=255"`
	CreatedAt   time.Time `validate:"structonly"`
	OwnerID     *int64    `validate:"omitempty,gt=0"`
}

// SendMessage represents model for message broker
type SendMessage struct {
	ChatID   int64  `validate:"required,gt=0"`
	SenderID int64  `validate:"required,gt=0"`
	Message  string `validate:"required,min=1,max=4096"`
}
