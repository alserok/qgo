package qgo

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func newRedisConsumer(addr, channel string, customs ...Customizer[any]) *redisConsumer {
	rc := redisConsumer{}
	for _, customizer := range customs {
		customizer(&rc)
	}

	cl := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: rc.password,
		DB:       rc.db,
	})
	rc.cl = cl

	if err := cl.Ping(context.Background()).Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	sub := cl.Subscribe(context.Background(), channel)
	rc.sub = sub

	msgs := rc.sub.Channel()
	rc.msgs = msgs

	return &rc
}

type redisConsumer struct {
	cl   *redis.Client
	sub  *redis.PubSub
	msgs <-chan *redis.Message

	password string
	db       int
}

func (r *redisConsumer) Consume(ctx context.Context) (*Message, error) {
	select {
	case msg := <-r.msgs:
		return &Message{
			Body: []byte(msg.Payload),
		}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func (r *redisConsumer) Close() error {
	if err := r.cl.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	if err := r.sub.Close(); err != nil {
		return fmt.Errorf("failed  to close subscription: %w", err)
	}

	return nil
}
