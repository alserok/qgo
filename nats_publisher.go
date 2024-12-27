package qgo

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

func newNatsPublisher(addr, topic string, customs ...Customizer[any]) *natsPublisher {
	pub := natsPublisher{
		topic: topic,
	}
	for _, custom := range customs {
		custom(&pub)
	}

	if pub.retryAttempts != nil {
		pub.opts = append(pub.opts, nats.RetryAttempts(*pub.retryAttempts))
	}
	if pub.retryWait != nil {
		pub.opts = append(pub.opts, nats.RetryWait(*pub.retryWait))
	}

	nc, err := nats.Connect(addr)
	if err != nil {
		panic("failed to connect to NATS: " + err.Error())
	}
	pub.nc = nc

	js, err := nc.JetStream()
	if err != nil {
		panic("failed to connect to NATS JetStream: " + err.Error())
	}
	pub.js = js

	_, err = js.StreamInfo(topic)
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     pub.topic,
			Subjects: pub.subjects,
		})
		if err != nil {
			panic("failed to add stream: " + err.Error())
		}
	}

	return &pub
}

type natsPublisher struct {
	js nats.JetStreamContext
	nc *nats.Conn

	topic          string
	defaultSubject string
	subjects       []string

	retryWait     *time.Duration
	retryAttempts *int

	opts []nats.PubOpt
}

func (n *natsPublisher) Produce(ctx context.Context, message *Message) error {
	if message.subject == "" {
		message.subject = n.defaultSubject
	} else {
		message.subject = n.topic + "." + message.subject
	}

	_, err := n.js.Publish(message.subject, message.Body, n.opts...)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (n *natsPublisher) Close() error {
	n.js.CleanupPublisher()
	n.nc.Close()
	return nil
}
