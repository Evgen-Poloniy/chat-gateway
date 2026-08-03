package kafka

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func NewKafkaProducer(cfg *Config, opts ...Option) (*kafka.Writer, error) {
	if cfg == nil {
		return nil, errors.New("provided kafka config is nil")
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("failed to apply security option: %w", err)
		}
	}

	var mechanism sasl.Mechanism
	var err error

	switch cfg.SASLMechanism {
	case "PLAIN":
		mechanism = plain.Mechanism{
			Username: cfg.SASLUsername,
			Password: cfg.SASLPassword,
		}
	case "SCRAM-SHA-256":
		mechanism, err = scram.Mechanism(scram.SHA256, cfg.SASLUsername, cfg.SASLPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to create scram-sha-256 mechanism: %w", err)
		}
	case "SCRAM-SHA-512":
		mechanism, err = scram.Mechanism(scram.SHA512, cfg.SASLUsername, cfg.SASLPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to create scram-sha-512 mechanism: %w", err)
		}
	}

	transport := &kafka.Transport{
		SASL: mechanism,
	}

	if cfg.SecurityProtocol == "SASL_SSL" || cfg.SecurityProtocol == "SSL" {
		transport.TLS = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.BootstrapServers...),
		MaxAttempts:  cfg.Retries,
		BatchSize:    cfg.BatchNumMessages,
		BatchTimeout: time.Duration(cfg.LingerMs) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.MessageTimeoutMs) * time.Millisecond,
		RequiredAcks: parseAcks(cfg.Acks),
		Compression:  parseCompression(cfg.CompressionType),
		Transport:    transport,
		Balancer:     &kafka.LeastBytes{},
	}

	return writer, nil
}

func parseAcks(acks string) kafka.RequiredAcks {
	switch acks {
	case "all", "-1":
		return kafka.RequireAll
	case "1":
		return kafka.RequireOne
	case "0":
		return kafka.RequireNone
	default:
		return kafka.RequireAll
	}
}

func parseCompression(codec string) kafka.Compression {
	switch codec {
	case "lz4":
		return kafka.Lz4
	case "gzip":
		return kafka.Gzip
	case "snappy":
		return kafka.Snappy
	case "zstd":
		return kafka.Zstd
	default:
		return 0
	}
}
