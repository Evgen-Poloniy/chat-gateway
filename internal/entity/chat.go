package entity

import "time"

// DirectChat represents database model for creating direct chat.
type DirectChat struct {
	ChatID      int64
	SenderID    int64
	RecipientID int64
	CreatedAt   time.Time
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
