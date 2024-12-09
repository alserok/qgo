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

	js, err := nc.JetStream()
	if err != nil {
		panic("failed to connect to JetStream: " + err.Error())
	}

	chMessages := make(chan *Message, 100)
	chErrors := make(chan error, 1)

	cons.chMessages = chMessages
	cons.chErrors = chErrors

	cc, err := js.Subscribe(cons.subject, func(msg *nats.Msg) {
		chMessages <- &Message{
			Body: msg.Data,
		}
		if err = msg.Ack(); err != nil {
			chErrors <- err
		}
	}, nats.ManualAck())
	if err != nil {
		panic("failed to consume: " + err.Error())
	}
	cons.cc = cc

	return &cons
}

type natsConsumer struct {
	cc *nats.Subscription
	nc *nats.Conn

	chMessages chan *Message
	chErrors   chan error

	topic   string
	subject string
}

func (n *natsConsumer) Consume(ctx context.Context) (*Message, error) {
	select {
	case msg, ok := <-n.chMessages:
		if !ok {
			return nil, fmt.Errorf("consumer closed")
		}
		return msg, nil
	case err := <-n.chErrors:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func (n *natsConsumer) Close() error {
	n.nc.Close()
	if err := n.cc.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	close(n.chErrors)
	close(n.chMessages)

	return nil
}
