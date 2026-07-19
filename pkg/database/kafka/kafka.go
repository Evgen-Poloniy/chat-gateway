package kafka

import (
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewKafkaProducer(configMap *kafka.ConfigMap, opts ...Option) (*kafka.Producer, error) {
	if configMap == nil {
		return nil, errors.New("provided kafka config map is nil")
	}

	for _, opt := range opts {
		if err := opt(configMap); err != nil {
			return nil, fmt.Errorf("failed to apply security option: %w", err)
		}
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return producer, nil
}

func GetTopicName(topic *string) string {
	if topic == nil {
		return "unknown"
	}
	return *topic
}
