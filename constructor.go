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
	Body      []byte
	ID        string
	Timestamp time.Time
	Partition int32
}

func (m *Message) DecodeBody(target any) error {
	if err := json.Unmarshal(m.Body, target); err != nil {
		return fmt.Errorf("failed to decode message body: %w", err)
	}

	return nil
}

type Producer interface {
	Produce(ctx context.Context, message *Message) error
	Close() error
}

func NewProducer(t uint, addr, topic string, customs ...Customizer[any]) Producer {
	switch t {
	case Kafka:
		return newKafkaProducer(addr, topic)
	case Nats:
		return newNatsPublisher(addr, topic, customs...)
	case Redis:
		return nil
	case Rabbit:
		return nil
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
		return nil
	case Rabbit:
		return nil
	default:
		panic("invalid Consumer type")
	}
}

type Customizer[T any] func(T)
