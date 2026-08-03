package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/segmentio/kafka-go"
)

// SendMessage sends message into target chat.
func (k *KafkaRepository) SendMessage(ctx context.Context, message *model.SendMessage) error {
	values, err := json.Marshal(message)
	if err != nil {
		return errs.NewAppError(
			errs.CodeSerializationError,
			"message broker error: failed to serialize incoming message data",
			fmt.Errorf("message broker error: %w", err),
		)
	}

	kafkaMsg := kafka.Message{
		Topic: "messenger-chat",
		Key:   message.ChatID[:],
		Value: values,
	}

	if err := k.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return errs.NewAppError(
			errs.CodeKafkaError,
			"message broker error: failed to send message",
			fmt.Errorf("message broker error: %w", err),
		)
	}

	return nil
}
