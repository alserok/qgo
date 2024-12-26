package qgo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	Kafka = iota
	Rabbit
	Redis
	Nats
)

type Message struct {
	Body      []byte    `json:"body"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`

	// Kafka message properties
	partition int32

	// Rabbit message properties
	ack      func() error
	exchange string

	// Redis message properties
	channel string

	// Nats message properties
	subject string
}

func (m *Message) SetSubject(sub string) {
	m.subject = sub
}

func (m *Message) Ack() error {
	return m.ack()
}

func (m *Message) SetExchange(exchange string) {
	m.exchange = exchange
}

func (m *Message) SetChannel(ch string) {
	m.channel = ch
}

func (m *Message) SetPartition(part int32) {
	m.partition = part
}

func (m *Message) DecodeBody(target any) error {
	if err := json.Unmarshal(m.Body, target); err != nil {
		return fmt.Errorf("failed to decode message body: %w", err)
	}

	return nil
}

func (m *Message) EncodeToBody(in any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to encode to body: %w", err)
	}

	m.Body = b

	return nil
}

type Producer interface {
	Produce(ctx context.Context, message *Message) error
	Close() error
}

func NewProducer(t uint, addr, topic string, customs ...Customizer[any]) Producer {
	switch t {
	case Kafka:
		return newKafkaProducer(addr, topic, customs...)
	case Nats:
		return newNatsPublisher(addr, topic, customs...)
	case Redis:
		return newRedisPublisher(addr, topic, customs...)
	case Rabbit:
		return newRabbitPublisher(addr, topic, customs...)
	default:
		panic("invalid Producer type")
	}
}

type Consumer interface {
	Consume(ctx context.Context) (*Message, error)
	Close() error
}

func NewConsumer(t uint, addr, topic string, customs ...Customizer[any]) Consumer {
	switch t {
	case Kafka:
		return newKafkaConsumer(addr, topic, customs...)
	case Nats:
		return newNatsConsumer(addr, topic, customs...)
	case Redis:
		return newRedisConsumer(addr, topic)
	case Rabbit:
		return newRabbitConsumer(addr, topic, customs...)
	default:
		panic("invalid Consumer type")
	}
}

type Customizer[T any] func(T)
