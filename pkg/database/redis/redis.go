package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// New initializes and returns a new Redis client.
func NewRedisCache(opts ...Option) (*redis.Client, error) {
	options := &Options{
		addr:         "localhost:6379",
		dialTimeout:  1 * time.Second,
		readTimeout:  1 * time.Second,
		writeTimeout: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         options.addr,
		Password:     options.password,
		DB:           options.db,
		DialTimeout:  options.dialTimeout,
		ReadTimeout:  options.readTimeout,
		WriteTimeout: options.writeTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), options.dialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}
