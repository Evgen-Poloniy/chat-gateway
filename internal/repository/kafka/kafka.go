package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// PostgresRepository represents implementation of interface MessengerRepository
type KafkaRepository struct {
	producer *kafka.Producer
}

func NewKafkaRepository(producer *kafka.Producer) *KafkaRepository {
	return &KafkaRepository{
		producer: producer,
	}
}
