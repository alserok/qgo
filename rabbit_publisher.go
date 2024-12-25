package qgo

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func newRabbitPublisher(addr, topic string, customs ...Customizer[any]) *rabbitPublisher {
	rp := rabbitPublisher{}
	for _, custom := range customs {
		custom(&rp)
	}

	conn, err := amqp.Dial(addr)
	if err != nil {
		panic("failed to dial: " + err.Error())
	}
	rp.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		panic("failed to init channel: " + err.Error())
	}
	rp.ch = ch

	q, err := ch.QueueDeclare(topic, rp.durable, rp.autoDelete, rp.exclusive, rp.noWait, nil)
	if err != nil {
		panic("failed to init queue: " + err.Error())
	}
	rp.q = q

	return &rp
}

type rabbitPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	exchange string

	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool

	mandatory bool
	immediate bool
}

func (r *rabbitPublisher) Produce(ctx context.Context, message *Message) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context canceled: %w", ctx.Err())
	}

	err := r.ch.PublishWithContext(ctx, r.exchange, r.q.Name, r.mandatory, r.immediate, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message.Body,
		Timestamp:   message.Timestamp,
		MessageId:   message.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *rabbitPublisher) Close() error {
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close amqp connection: %w", err)
	}

	if err := r.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	return nil
}
