package qgo

import (
	"context"
	"encoding/json"
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
	c         KafkaCustomizer
}

func (s *KafkaSuite) SetupTest() {
	s.addr, s.container = s.setupKafka()
	s.c = NewKafkaCustomizer()
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

func (s *KafkaSuite) TestWithCustomizers() {
	p := NewProducer(Kafka, s.addr, "test", s.c.WithFlushFrequency(400*time.Millisecond), s.c.WithCompression(CompressionGZIP), s.c.WithRequiredAcks(AckWaitForAll))
	c := NewConsumer(Kafka, s.addr, "test", s.c.WithOffset(OffsetOldest), s.c.WithPartition(0))
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

func (s *KafkaSuite) TestMessageDecode() {
	p := NewProducer(Kafka, s.addr, "test")
	c := NewConsumer(Kafka, s.addr, "test")
	defer func() {
		s.Require().NoError(p.Close())
		s.Require().NoError(c.Close())
	}()

	messagesAmount := 5
	for i := range messagesAmount {
		type body struct {
			String string
			Int    int
			Bool   bool
			Str    struct {
				Field string
			}
		}

		data := body{
			String: fmt.Sprintf("number: %d", i),
			Int:    i,
			Bool:   i%2 == 0,
			Str: struct{ Field string }{
				Field: fmt.Sprintf("Struct field: %d", i),
			},
		}

		b, err := json.Marshal(data)
		s.Require().NoError(err)

		s.Require().NoError(p.Produce(context.Background(), &Message{
			Body:      b,
			ID:        fmt.Sprintf("%d", i),
			Timestamp: time.Now(),
		}))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)

		var resData body
		s.Require().NoError(msg.DecodeBody(&resData))
		s.Require().Equal(data, resData)
	}
}

func (s *KafkaSuite) TestMessageEncode() {
	p := NewProducer(Kafka, s.addr, "test")
	c := NewConsumer(Kafka, s.addr, "test")
	defer func() {
		s.Require().NoError(p.Close())
		s.Require().NoError(c.Close())
	}()

	messagesAmount := 5
	for i := range messagesAmount {
		type body struct {
			String string
			Int    int
			Bool   bool
			Str    struct {
				Field string
			}
		}

		data := body{
			String: fmt.Sprintf("number: %d", i),
			Int:    i,
			Bool:   i%2 == 0,
			Str: struct{ Field string }{
				Field: fmt.Sprintf("Struct field: %d", i),
			},
		}

		msg := &Message{
			ID:        fmt.Sprintf("%d", i),
			Timestamp: time.Now(),
		}
		s.Require().NoError(msg.EncodeToBody(data))

		s.Require().NoError(p.Produce(context.Background(), msg))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)

		var resData body
		s.Require().NoError(json.Unmarshal(msg.Body, &resData))
		s.Require().Equal(data, resData)
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
