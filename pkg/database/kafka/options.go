package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Option defines a function type to safely mutate the base kafka.ConfigMap.
type Option func(cm *kafka.ConfigMap) error

// WithSecurityProtocol sets the communication protocol (e.g., SASL_SSL, SASL_PLAINTEXT).
// It ignores empty strings or standard unencrypted PLAINTEXT.
func WithSecurityProtocol(protocol string) Option {
	return func(cm *kafka.ConfigMap) error {
		if protocol == "" || protocol == "PLAINTEXT" {
			return nil
		}
		if err := cm.SetKey("security.protocol", protocol); err != nil {
			return fmt.Errorf("failed to set security.protocol: %w", err)
		}
		return nil
	}
}

// WithSASLMechanism sets the SASL authentication mechanism (e.g., PLAIN, SCRAM-SHA-512).
func WithSASLMechanism(mechanism string) Option {
	return func(cm *kafka.ConfigMap) error {
		if mechanism == "" {
			return nil
		}
		if err := cm.SetKey("sasl.mechanism", mechanism); err != nil {
			return fmt.Errorf("failed to set sasl.mechanism: %w", err)
		}
		return nil
	}
}

// WithSASLUsername sets the username for SASL authentication.
func WithSASLUsername(username string) Option {
	return func(cm *kafka.ConfigMap) error {
		if username == "" {
			return nil
		}
		if err := cm.SetKey("sasl.username", username); err != nil {
			return fmt.Errorf("failed to set sasl.username: %w", err)
		}
		return nil
	}
}

// WithSASLPassword sets the password for SASL authentication.
func WithSASLPassword(password string) Option {
	return func(cm *kafka.ConfigMap) error {
		if password == "" {
			return nil
		}
		if err := cm.SetKey("sasl.password", password); err != nil {
			return fmt.Errorf("failed to set sasl.password: %w", err)
		}
		return nil
	}
}
