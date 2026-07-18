package resolver

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/resolver_mocks.go -package=mock_resolver_repository

// Channel represents interface for work with broker channel.
type Channel interface {
	// SubscribeOnEventChannel subscribes on broker channel once fro all time of work application.
	SubscribeOnEventChannel(ctx context.Context)

	// ResolveEvent returns event from channel.
	ResolveEvent(ctx context.Context) (*model.EventMessage, error)
}

// Cache represents interface for work with cache with online users.
type Cache interface {
	// ExpireChatID expire key TTL.
	ExpireChatID(ctx context.Context, chatID string) error

	// GetChatMembers allows get user_ids from the messenger cache by chat_id.
	GetChatMembers(ctx context.Context, chatID string) ([]string, error)
}

// ResolverService represents implementation of Resolver interface.
type ResolverService struct {
	channel Channel
	cache   Cache
}

func NewResolverService(channel Channel, cache Cache) *ResolverService {
	return &ResolverService{
		channel: channel,
		cache:   cache,
	}
}
