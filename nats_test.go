package qgo

import (
	"context"
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
}

func (s *NatsSuite) SetupTest() {
	s.addr, s.container = s.setupNATS()
}

func (s *NatsSuite) TeardownTest() {
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func (s *NatsSuite) TestDefault() {
	p := NewProducer(Nats, s.addr, "TEST", WithSubject("a"))
	c := NewConsumer(Nats, s.addr, "TEST", WithSubject("a"))
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
	c := NewConsumer(Nats, s.addr, "TEST", WithSubject("a"))
	p := NewProducer(Nats, s.addr, "TEST", WithSubject("a"), WithRetryWait(time.Second), WithRetryAttempts(3))
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
