package database

import (
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	"github.com/google/uuid"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	MessageTopic = "message-topic"
)

func NewKafkaProducer(config *config.KafkaConfig, logger *logrus.Logger) (*kafka.Producer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":                     config.BootstrapServers,
		"acks":                                  config.Acks,
		"enable.idempotence":                    config.EnableIdempotence,
		"retries":                               config.Retries,
		"max.in.flight.requests.per.connection": config.MaxInFlightRequestsPerConn,
		"linger.ms":                             config.LingerMs,
		"batch.num.messages":                    config.BatchNumMessages,
		"compression.type":                      config.CompressionType,
		"queue.buffering.max.messages":          config.QueueBufferingMaxMessages,
		"message.timeout.ms":                    config.MessageTimeout,
		"num.partitions":                        config.NumPartitions,
	}

	if config.SecurityProtocol != "" && config.SecurityProtocol != "PLAINTEXT" {
		configMap.SetKey("security.protocol", config.SecurityProtocol)
		configMap.SetKey("sasl.mechanism", config.SASLMechanism)
		configMap.SetKey("sasl.username", config.User)
		configMap.SetKey("sasl.password", config.Password)
	}

	p, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.WithFields(map[string]interface{}{
						"id":        uuid.NewString(),
						"topic":     getTopicName(ev.TopicPartition.Topic),
						"partition": ev.TopicPartition.Partition,
						"code":      "kafka_error",
					}).Error(fmt.Sprintf("kafka delivery error: %v", ev.TopicPartition.Error))
				} else {
					logger.WithFields(map[string]interface{}{
						"id":        uuid.NewString(),
						"topic":     getTopicName(ev.TopicPartition.Topic),
						"partition": ev.TopicPartition.Partition,
						"offset":    ev.TopicPartition.Offset,
					}).Info("kafka message delivered")
				}
			}
		}
	}()

	return p, nil
}

func getTopicName(topic *string) string {
	if topic == nil {
		return "unknown"
	}
	return *topic
}
