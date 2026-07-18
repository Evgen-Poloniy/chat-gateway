package rds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
)

// SubscribeOnEventChannel subscribes on broker channel once fro all time of work application.
func (r *RedisCache) SubscribeOnEventChannel(ctx context.Context) {
	r.eventChan = make(chan model.Event, r.chatConf.ChatTtl)
	pubsub := r.rdb.Subscribe(ctx, "events:users:dispatch")

	go func() {
		defer func() {
			pubsub.Close()
			close(r.eventChan)
		}()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case redisMsg, ok := <-ch:
				if !ok {
					select {
					case <-ctx.Done():
					case r.eventChan <- model.Event{
						Err: errs.NewAppError(
							errs.CodeFailedEventChannel,
							errs.ErrFailedEventChannel.Error(),
							errs.ErrFailedEventChannel,
						),
					}:
					}
					return
				}

				var msg model.EventMessage
				if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
					select {
					case <-ctx.Done():
						return
					case r.eventChan <- model.Event{
						Err: errs.NewAppError(
							errs.CodeDeserializationError,
							"failed to unmarshal redis pubsub message",
							fmt.Errorf("failed to unmarshal redis pubsub message: %w", err),
						),
					}:
					}
					continue
				}

				select {
				case <-ctx.Done():
					return
				case r.eventChan <- model.Event{
					Data: msg,
				}:
				}
			}
		}
	}()
}

// ResolveEvent returns event from channel.
func (r *RedisCache) ResolveEvent(ctx context.Context) (*model.EventMessage, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ev, ok := <-r.eventChan:
		if !ok {
			return nil, errs.NewAppError(
				errs.CodeFailedEventChannel,
				errs.ErrFailedEventChannel.Error(),
				errs.ErrFailedEventChannel,
			)
		}

		if ev.Err != nil {
			return nil, ev.Err
		}

		return &ev.Data, nil
	}
}
