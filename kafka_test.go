package qgo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func TestKafka(t *testing.T) {
	suite.Run(t, new(KafkaSuite))
}

type KafkaSuite struct {
	suite.Suite

	addr      string
	container *kafka.KafkaContainer
}

func (s *KafkaSuite) SetupTest() {
	s.addr, s.container = s.setupKafka()
}

func (s *KafkaSuite) TeardownTest() {
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func (s *KafkaSuite) TestDefault() {
	p := NewProducer(Kafka, s.addr, "test")
	c := NewConsumer(Kafka, s.addr, "test")
	defer func() {
		s.Require().NoError(p.Close())
		s.Require().NoError(c.Close())
	}()

	messagesAmount := 5
	for i := range messagesAmount {
		s.Require().NoError(p.Produce(context.Background(), &Message{
			Body:      []byte("body"),
			ID:        fmt.Sprintf("%d", i),
			Timestamp: time.Now(),
		}))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)
	}
}

func (s *KafkaSuite) TestDefaultWithConsumerCustomizers() {
	p := NewProducer(Kafka, s.addr, "test", WithFlushFrequency(400*time.Millisecond), WithCompression(CompressionGZIP), WithRequiredAcks(AckWaitForAll))
	c := NewConsumer(Kafka, s.addr, "test", WithOffset(OffsetOldest), WithPartition(0))
	defer func() {
		s.Require().NoError(p.Close())
		s.Require().NoError(c.Close())
	}()

	messagesAmount := 5
	for i := range messagesAmount {
		s.Require().NoError(p.Produce(context.Background(), &Message{
			Body:      []byte("body"),
			ID:        fmt.Sprintf("%d", i),
			Timestamp: time.Now(),
		}))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)
	}
}

func (s *KafkaSuite) setupKafka() (string, *kafka.KafkaContainer) {
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
