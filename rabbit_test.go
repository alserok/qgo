package qgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func TestRabbit(t *testing.T) {
	suite.Run(t, new(RabbitSuite))
}

type RabbitSuite struct {
	suite.Suite

	addr      string
	container *rabbitmq.RabbitMQContainer
}

func (s *RabbitSuite) SetupTest() {
	s.addr, s.container = s.setupRabbit()
}

func (s *RabbitSuite) TeardownTest() {
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func (s *RabbitSuite) TestDefault() {
	p := NewProducer(Rabbit, s.addr, "test")
	c := NewConsumer(Rabbit, s.addr, "test")
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

func (s *RabbitSuite) TestWithCustomizers() {
	p := NewProducer(Rabbit, s.addr, "test")
	c := NewConsumer(Rabbit, s.addr, "test", WithAutoAcknowledgement())
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

func (s *RabbitSuite) TestMessageAck() {
	p := NewProducer(Rabbit, s.addr, "test")
	c := NewConsumer(Rabbit, s.addr, "test")
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
		s.Require().NoError(msg.Ack())
	}
}

func (s *RabbitSuite) TestMessageDecode() {
	p := NewProducer(Rabbit, s.addr, "test")
	c := NewConsumer(Rabbit, s.addr, "test", WithAutoAcknowledgement())
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

func (s *RabbitSuite) TestMessageEncode() {
	p := NewProducer(Rabbit, s.addr, "test")
	c := NewConsumer(Rabbit, s.addr, "test")
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

func (s *RabbitSuite) setupRabbit() (string, *rabbitmq.RabbitMQContainer) {
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
