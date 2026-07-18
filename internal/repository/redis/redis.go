package rds

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/redis/go-redis/v9"
)

// RedisCache represents implementation of interface MessengerCache
type RedisCache struct {
	rdb            *redis.Client
	chatConf       *config.RedisConfig
	resolverConfig *config.ResolverConfig
	eventChan      chan model.Event
}

func NewRedisCache(rdb *redis.Client, chatConf *config.RedisConfig, resolverConfig *config.ResolverConfig) *RedisCache {
	return &RedisCache{
		rdb:            rdb,
		chatConf:       chatConf,
		resolverConfig: resolverConfig,
	}
}
