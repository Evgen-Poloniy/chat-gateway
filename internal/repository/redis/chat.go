package rds

import (
	"context"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/pkg/errs"
)

// AddChatMembers writes user_ids at cache by key chat_id.
func (r *RedisCache) AddChatMembers(ctx context.Context, chatID string, userIDs []string) error {
	key := fmt.Sprintf("chat:%s:members", chatID)

	pipe := r.rdb.Pipeline()

	pipe.SAdd(ctx, key, userIDs)

	pipe.Expire(ctx, key, r.chatConf.ChatTtl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errs.NewAppError(
			errs.CodeRedisError,
			"cache error: failed to add user",
			fmt.Errorf("cache error: failed to add user: %w", err),
		)
	}
	return nil
}

// RemoveChatUser removes user from Redis by chat_id.
func (r *RedisCache) RemoveChatUser(ctx context.Context, chatID string, userID string) error {
	key := fmt.Sprintf("chat:%s:members", chatID)

	pipe := r.rdb.Pipeline()

	pipe.SRem(ctx, key, userID)

	pipe.Expire(ctx, key, r.chatConf.ChatTtl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errs.NewAppError(
			errs.CodeRedisError,
			"cache error: failed to remove user",
			fmt.Errorf("cache error: failed to remove user: %w", err),
		)
	}

	return nil
}

// ExpireChatID expire key TTL.
func (r *RedisCache) ExpireChatID(ctx context.Context, chatID string) error {
	key := fmt.Sprintf("chat:%s:members", chatID)

	_, err := r.rdb.Expire(ctx, key, r.chatConf.ChatTtl).Result()
	if err != nil {
		return errs.NewAppError(
			errs.CodeRedisError,
			"cache error: failed to expire key",
			fmt.Errorf("cache error: failed to expire key: %w", err),
		)
	}

	return nil
}

// GetChatMembers allows get user_ids from the messenger cache by chat_id.
func (r *RedisCache) GetChatMembers(ctx context.Context, chatID string) ([]string, error) {
	key := fmt.Sprintf("chat:%s:members", chatID)

	userIDs, err := r.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errs.NewAppError(
			errs.CodeRedisError,
			"cache error: failed to expire key",
			fmt.Errorf("cache error: failed to get chat users: %w", err),
		)
	}

	return userIDs, nil
}
