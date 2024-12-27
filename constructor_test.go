package qgo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"testing"
)

func TestConstructor(t *testing.T) {
	suite.Run(t, new(ConstructorSuite))
}

type ConstructorSuite struct {
	suite.Suite
}

func (s *ConstructorSuite) TestKafkaConstructor() {
	addr, container := s.setupKafka()
	defer func() {
		s.Require().NoError(container.Terminate(context.Background()))
	}()

	producer := NewProducer(Kafka, addr, "topic")
	s.Require().NotNil(producer)
	consumer := NewConsumer(Kafka, addr, "topic")
	s.Require().NotNil(consumer)
}

func (s *ConstructorSuite) TestRabbitConstructor() {
	addr, container := s.setupRabbit()
	defer func() {
		s.Require().NoError(container.Terminate(context.Background()))
	}()

	producer := NewProducer(Rabbit, addr, "topic")
	s.Require().NotNil(producer)
	consumer := NewConsumer(Rabbit, addr, "topic")
	s.Require().NotNil(consumer)
}

func (s *ConstructorSuite) TestNatsConstructor() {
	addr, container := s.setupNATS()
	defer func() {
		s.Require().NoError(container.Terminate(context.Background()))
	}()

	c := NewNATSCustomizer()

	producer := NewProducer(Nats, addr, "topic", c.WithSubjects("subj"))
	s.Require().NotNil(producer)
	consumer := NewConsumer(Nats, addr, "topic", c.WithSubjects("subj"))
	s.Require().NotNil(consumer)
}

func (s *ConstructorSuite) TestRedisConstructor() {
	addr, container := s.setupRedis()
	defer func() {
		s.Require().NoError(container.Terminate(context.Background()))
	}()

	producer := NewProducer(Redis, addr, "topic")
	s.Require().NotNil(producer)
	consumer := NewConsumer(Redis, addr, "topic")
	s.Require().NotNil(consumer)
}

func (s *ConstructorSuite) setupRedis() (string, *redis.RedisContainer) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx,
		"redis:latest",
	)
	s.Require().NoError(err)
	s.Require().NotNil(redisContainer)
	s.Require().True(redisContainer.IsRunning())

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	s.Require().NoError(err)

	return fmt.Sprintf("localhost:%s", port.Port()), redisContainer
}

func (s *ConstructorSuite) setupKafka() (string, *kafka.KafkaContainer) {
	ctx := context.Background()

	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)
	s.Require().NoError(err)
	s.Require().NotNil(kafkaContainer)
	s.Require().True(kafkaContainer.IsRunning())

	port, err := kafkaContainer.MappedPort(ctx, "9093/tcp")
	s.Require().NoError(err)

	return fmt.Sprintf("localhost:%s", port.Port()), kafkaContainer
}

func (s *ConstructorSuite) setupNATS() (string, *nats.NATSContainer) {
	ctx := context.Background()

	natsContainer, err := nats.Run(ctx,
		"nats:2.9",
	)
	s.Require().NoError(err)
	s.Require().NotNil(natsContainer)
	s.Require().True(natsContainer.IsRunning())

	port, err := natsContainer.MappedPort(ctx, "4222/tcp")
	s.Require().NoError(err)

	return fmt.Sprintf("localhost:%s", port.Port()), natsContainer
}

func (s *ConstructorSuite) setupRabbit() (string, *rabbitmq.RabbitMQContainer) {
	ctx := context.Background()

	rabbitContainer, err := rabbitmq.Run(ctx,
		"rabbitmq:latest",
	)
	s.Require().NoError(err)
	s.Require().NotNil(rabbitContainer)
	s.Require().True(rabbitContainer.IsRunning())

	port, err := rabbitContainer.MappedPort(ctx, "5672/tcp")
	s.Require().NoError(err)

	return fmt.Sprintf("amqp://guest:guest@localhost:%s/", port.Port()), rabbitContainer
}
