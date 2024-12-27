package qgo

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
)

func newNatsConsumer(addr, topic string, customs ...Customizer[any]) *natsConsumer {
	cons := natsConsumer{
		topic: topic,
	}
	for _, custom := range customs {
		custom(&cons)
	}

	nc, err := nats.Connect(addr)
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	chMessages := make(chan *nats.Msg, 128)

	sub, err := nc.ChanSubscribe(cons.subject, chMessages)
	if err != nil {
		panic("failed to subscribe to channel: " + err.Error())
	}
	cons.cc = sub
	cons.ch = chMessages

	return &cons
}

type natsConsumer struct {
	cc *nats.Subscription
	nc *nats.Conn

	ch chan *nats.Msg

	topic   string
	subject string
}

func (n *natsConsumer) Consume(ctx context.Context) (*Message, error) {
	select {
	case msg, ok := <-n.ch:
		if !ok {
			return nil, fmt.Errorf("consumer closed")
		}
		return &Message{
			Body: msg.Data,
		}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func (n *natsConsumer) Close() error {
	n.nc.Close()
	if err := n.cc.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	close(n.ch)

	return nil
}
