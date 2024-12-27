package qgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/nats"
)

func TestNats(t *testing.T) {
	suite.Run(t, new(NatsSuite))
}

type NatsSuite struct {
	suite.Suite

	addr      string
	container *nats.NATSContainer
	c         NatsCustomizer
}

func (s *NatsSuite) SetupTest() {
	s.addr, s.container = s.setupNATS()
	s.c = NewNATSCustomizer()
}

func (s *NatsSuite) TeardownTest() {
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func (s *NatsSuite) TestDefault() {
	c := NewConsumer(Nats, s.addr, "TEST", s.c.WithSubjects("a"))
	p := NewProducer(Nats, s.addr, "TEST", s.c.WithSubjects("a"))
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

func (s *NatsSuite) TestWithCustomizers() {
	c := NewConsumer(Nats, s.addr, "TEST", s.c.WithSubjects("a"))
	p := NewProducer(Nats, s.addr, "TEST", s.c.WithSubjects("a", "b"), s.c.WithRetryWait(time.Second), s.c.WithRetryAttempts(3))
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

func (s *NatsSuite) TestWithMessageSubject() {
	c := NewConsumer(Nats, s.addr, "test", s.c.WithSubjects("b"))
	p := NewProducer(Nats, s.addr, "test", s.c.WithSubjects("a", "b"))
	defer func() {
		s.Require().NoError(p.Close())
		s.Require().NoError(c.Close())
	}()

	messagesAmount := 5
	for i := range messagesAmount {
		msg := &Message{
			Body:      []byte("body"),
			ID:        fmt.Sprintf("%d", i),
			Timestamp: time.Now(),
		}
		msg.SetSubject("b")

		s.Require().NoError(p.Produce(context.Background(), msg))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)
	}
}

func (s *NatsSuite) TestMessageDecode() {
	c := NewConsumer(Nats, s.addr, "test", s.c.WithSubjects("a"))
	p := NewProducer(Nats, s.addr, "test", s.c.WithSubjects("a"))
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

func (s *NatsSuite) TestMessageEncode() {
	c := NewConsumer(Nats, s.addr, "test", s.c.WithSubjects("a"))
	p := NewProducer(Nats, s.addr, "test", s.c.WithSubjects("a"))
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

func (s *NatsSuite) setupNATS() (string, *nats.NATSContainer) {
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
