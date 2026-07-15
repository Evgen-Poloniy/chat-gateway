package kf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/database"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// SendMessage sends message into target chat.
func (k *KafkaRepository) SendMessage(ctx context.Context, message *model.SendMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		values, err := json.Marshal(message)
		if err != nil {
			return errs.NewAppError(
				errs.CodeSerializationError,
				"message broker error: failed to serialize incoming message data",
				fmt.Errorf("message broker error: %v", err),
			)
		}

		topic := database.MessageTopic

		return k.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          values,
			Key:            message.ChatID[:],
		}, nil)
	}
}
