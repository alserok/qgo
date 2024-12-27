package qgo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func newRedisConsumer(addr, channel string, customs ...Customizer[any]) *redisConsumer {
	rc := redisConsumer{}
	for _, customizer := range customs {
		customizer(&rc)
	}

	cl := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   rc.password,
		DB:         rc.db,
		MaxRetries: rc.maxRetries,
		TLSConfig:  rc.tlsConfig,
		Network:    rc.network,
		ClientName: rc.clientName,
	})
	rc.cl = cl

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := cl.Ping(ctx).Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	sub := cl.Subscribe(ctx, channel)
	rc.sub = sub

	msgs := rc.sub.Channel()
	rc.msgs = msgs

	return &rc
}

type redisConsumer struct {
	cl   *redis.Client
	sub  *redis.PubSub
	msgs <-chan *redis.Message

	password   string
	db         int
	maxRetries int
	tlsConfig  *tls.Config
	network    string
	clientName string
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
