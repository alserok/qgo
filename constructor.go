package qgo

import (
	"context"
	"crypto/tls"
	"time"
)

const (
	Kafka = iota
	Rabbit
	Redis
	Nats
)

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

type KafkaCustomizer interface {
	WithOffset(offset int64) Customizer[any]
	WithPartition(partition int32) Customizer[any]
	WithRequiredAcks(acks int) Customizer[any]
	WithCompression(comp int) Customizer[any]
	WithFlushFrequency(freq time.Duration) Customizer[any]
	WithIdempotent() Customizer[any]
}

func NewKafkaCustomizer() KafkaCustomizer {
	return kafkaCustomizer{}
}

type RabbitCustomizer interface {
	WithNoWait() Customizer[any]
	WithExclusive() Customizer[any]
	WithAutoDelete() Customizer[any]
	WithDurable() Customizer[any]
	WithExchange(exchange string) Customizer[any]
	WithMandatory() Customizer[any]
	WithImmediate() Customizer[any]
	WithConsumerTag(tag string) Customizer[any]
	WithNoLocal() Customizer[any]
	WithAutoAcknowledgement() Customizer[any]
	WithConsumerArguments(args map[string]interface{}) Customizer[any]
}

func NewRabbitCustomizer() RabbitCustomizer {
	return rabbitCustomizer{}
}

type NatsCustomizer interface {
	WithSubjects(subjects ...string) Customizer[any]
	WithRetryWait(wait time.Duration) Customizer[any]
	WithRetryAttempts(attempts int) Customizer[any]
}

func NewNATSCustomizer() NatsCustomizer {
	return natsCustomizer{}
}

type RedisCustomizer interface {
	WithPassword(password string) Customizer[any]
	WithDB(db int) Customizer[any]
	WithRetries(retries int) Customizer[any]
	WithTLSConfig(cfg *tls.Config) Customizer[any]
	WithNetwork(net string) Customizer[any]
	WithClientName(name string) Customizer[any]
}

func NewRedisCustomizer() RedisCustomizer {
	return redisCustomizer{}
}
