package qgo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func newRedisPublisher(addr, channel string, customs ...Customizer[any]) *redisPublisher {
	rp := redisPublisher{
		defaultChannel: channel,
	}
	for _, customizer := range customs {
		customizer(&rp)
	}

	cl := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   rp.password,
		DB:         rp.db,
		MaxRetries: rp.maxRetries,
		TLSConfig:  rp.tlsConfig,
		Network:    rp.network,
		ClientName: rp.clientName,
	})
	rp.cl = cl

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := cl.Ping(ctx).Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	return &rp
}

type redisPublisher struct {
	cl *redis.Client

	defaultChannel string
	password       string
	db             int
	maxRetries     int
	tlsConfig      *tls.Config
	network        string
	clientName     string
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
