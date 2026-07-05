package kf

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// PostgresRepository represents implementation of interface MessengerRepository
type KafkaRepository struct {
	p *kafka.Producer
}

func NewKafkaRepository(p *kafka.Producer) *KafkaRepository {
	return &KafkaRepository{
		p: p,
	}
}
