package entity

import (
	"fmt"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/google/uuid"
)

// EventMessage represents event message with message recipients
type EventMessage struct {
	MessageID string
	ChatID    string
	UserIDs   uuid.UUIDs
	SenderID  string
	Text      string
	CreatedAt time.Time
}

func (e *EventMessage) ParseUUIDs(strs []string) error {
	if len(strs) == 0 {
		return errs.NewAppError(
			errs.CodeEmptyUserIDs,
			"validation error: "+errs.ErrEmptyUserIDs.Error(),
			fmt.Errorf("validation error: %w", errs.ErrEmptyUserIDs),
		)
	}

	e.UserIDs = make(uuid.UUIDs, 0, len(strs))

	for _, userID := range strs {
		id, err := uuid.Parse(userID)
		if err != nil {
			continue
		}

		e.UserIDs = append(e.UserIDs, id)
	}

	if len(e.UserIDs) == 0 {
		return errs.NewAppError(
			errs.CodeValidationError,
			"validation error: "+errs.ErrUserIDsConversionFailure.Error(),
			fmt.Errorf("validation error: %w", errs.ErrUserIDsConversionFailure),
		)
	}

	return nil
}
