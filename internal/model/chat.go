package model

import "time"

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID         int64
	ParticipantIDs []int64
	CreatedAt      time.Time
}

// GroupChat represents database model for creating direct chat.
type GroupChat struct {
	ChatID         int64
	ParticipantIDs []int64
	Name           string
	Title          *string
	Description    *string
	CreatedAt      time.Time
	OwnerID        int64
}

// Chat represents model of chat
type Chat struct {
	ChatID      int64     `db:"id"`
	ChatType    string    `db:"type"`
	Name        *string   `db:"name"`
	Title       *string   `db:"title"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	OwnerID     *int64    `db:"owner_id"`
}

// UpdateGroupChat represents database model for update data about group chat.
type UpdateGroupChat struct {
	UserIDUpdater int64     `db:"user_id_updater"`
	ChatID        int64     `db:"id"`
	Name          *string   `db:"name"`
	Title         *string   `db:"title"`
	Description   *string   `db:"description"`
	CreatedAt     time.Time `db:"created_at"`
	OwnerID       *int64    `db:"owner_id"`
}

// SendMessage represents model for message broker
type SendMessage struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	Message  string `json:"message" binding:"required"`
}
