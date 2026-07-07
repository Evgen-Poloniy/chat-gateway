package kf

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/database"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// SendMessage sends message into target chat.
func (k *KafkaRepository) SendMessage(message *dto.SendMessageReq) error {
	values, err := json.Marshal(message)
	if err != nil {
		return errs.NewAppError(
			"SERIALIZATION_ERROR",
			"message broker error: failed to serialize incoming message data",
			fmt.Sprintf("message broker error: %v", err),
			errs.ErrSerialization,
		)
	}

	topic := database.MessageTopic

	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          values,
		Key:            []byte(strconv.FormatInt(message.ChatID, 10)),
	}, nil)
}
