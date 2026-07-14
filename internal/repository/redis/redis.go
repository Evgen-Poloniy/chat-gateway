package rds

import "github.com/redis/go-redis/v9"

// RedisCache represents implementation of interface MessengerCache
type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}
