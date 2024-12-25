package qgo

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func newRabbitConsumer(addr, topic string, customs ...Customizer[any]) *rabbitConsumer {
	rc := rabbitConsumer{}
	for _, custom := range customs {
		custom(&rc)
	}

	conn, err := amqp.Dial(addr)
	if err != nil {
		panic("failed to dial: " + err.Error())
	}
	rc.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		panic("failed to init channel: " + err.Error())
	}
	rc.ch = ch

	q, err := ch.QueueDeclare(topic, rc.durable, rc.autoDelete, rc.exclusive, rc.noWait, nil)
	if err != nil {
		panic("failed to init queue: " + err.Error())
	}
	rc.q = q

	msgs, err := ch.Consume(rc.q.Name, rc.tag, rc.autoAcknowledgement, rc.exclusive, rc.noLocal, rc.noWait, nil)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}
	rc.msgs = msgs

	return &rc
}

type rabbitConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
	msgs <-chan amqp.Delivery

	topic string
	tag   string

	durable             bool
	autoDelete          bool
	exclusive           bool
	noWait              bool
	noLocal             bool
	autoAcknowledgement bool
}

func (r *rabbitConsumer) Consume(ctx context.Context) (*Message, error) {
	select {
	case msg, ok := <-r.msgs:
		if !ok {
			return nil, fmt.Errorf("consumer closed")
		}

		return &Message{
			Body:      msg.Body,
			ID:        msg.MessageId,
			Timestamp: msg.Timestamp,
			Ack: func() error {
				if err := msg.Ack(false); err != nil {
					return fmt.Errorf("failed to acknowledge: %w", err)
				}

				return nil
			},
		}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func (r *rabbitConsumer) Close() error {
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close amqp connection: %w", err)
	}

	if err := r.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	return nil
}
