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
)

type Message struct {
	Body      []byte
	ID        string
	Timestamp time.Time
	Partition int32
}

func (m *Message) DecodeBody(target *any) error {
	if err := json.Unmarshal(m.Body, target); err != nil {
		return fmt.Errorf("failed to decode message body: %w", err)
	}

	return nil
}

type Producer interface {
	Produce(ctx context.Context, message *Message) error
	Close() error
}

func NewProducer(t uint, addr, topic string) Producer {
	switch t {
	case Kafka:
		return newKafkaProducer(addr, topic)
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

func NewConsumer(t uint, addr, topic string, cfg ...any) Consumer {
	switch t {
	case Kafka:
		return newKafkaConsumer(addr, topic, cfg...)
	case Redis:
		return nil
	case Rabbit:
		return nil
	default:
		panic("invalid Consumer type")
	}
}

type KafkaConsumerConfig struct {
	Partition int32
	Offset    int64
}

const (
	OffsetNewest = -1
	OffsetOldest = -2
)
