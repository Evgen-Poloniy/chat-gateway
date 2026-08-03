package kafka

import (
	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	writer *kafka.Writer
}

func NewKafkaRepository(writer *kafka.Writer) *KafkaRepository {
	return &KafkaRepository{
		writer: writer,
	}
}
