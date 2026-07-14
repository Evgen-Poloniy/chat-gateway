package database

import (
	"context"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Host + ":" + config.Port,
		Password:     config.Password,
		DB:           0,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}
