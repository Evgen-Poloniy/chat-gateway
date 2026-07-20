package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
)

// SubscribeToEventChannel subscribes on broker channel once fro all time of work application.
func (r *RedisCache) SubscribeToEventChannel(ctx context.Context) error {
	pubsub := r.rdb.Subscribe(ctx, "events:users:dispatch")

	if _, err := pubsub.Receive(ctx); err != nil {
		if errClose := pubsub.Close(); errClose != nil {
			return errors.Join(err, errClose)
		}
		return err
	}

	go func() {
		defer func() {
			_ = pubsub.Close()
			close(r.events)
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
					case r.events <- model.Event{
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
					case r.events <- model.Event{
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
				case r.events <- model.Event{
					Data: msg,
				}:
				}
			}
		}
	}()

	return nil
}

// ResolveEvent returns event from channel.
func (r *RedisCache) ResolveEvent(ctx context.Context) (*model.EventMessage, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ev, ok := <-r.events:
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
