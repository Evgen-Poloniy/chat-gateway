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
