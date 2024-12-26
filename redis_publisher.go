package qgo

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func newRedisPublisher(addr, channel string, customs ...Customizer[any]) *redisPublisher {
	rp := redisPublisher{
		defaultChannel: channel,
	}
	for _, customizer := range customs {
		customizer(&rp)
	}

	cl := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: rp.password,
		DB:       rp.db,
	})
	rp.cl = cl

	if err := cl.Ping(context.Background()).Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	return &rp
}

type redisPublisher struct {
	cl *redis.Client

	defaultChannel string
	password       string
	db             int
}

func (r *redisPublisher) Produce(ctx context.Context, message *Message) error {
	if message.channel == "" {
		message.channel = r.defaultChannel
	}

	if err := r.cl.Publish(ctx, message.channel, message.Body).Err(); err != nil {
		return fmt.Errorf("faield to publish message: %w", err)
	}

	return nil
}

func (r *redisPublisher) Close() error {
	if err := r.cl.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}
