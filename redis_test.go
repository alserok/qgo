package qgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"testing"
	"time"
)

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

type RedisSuite struct {
	suite.Suite

	addr      string
	container *redis.RedisContainer
}

func (s *RedisSuite) SetupTest() {
	s.addr, s.container = s.setupRedis()
}

func (s *RedisSuite) TeardownTest() {
	s.Require().NoError(s.container.Terminate(context.Background()))
}

func (s *RedisSuite) TestDefault() {
	p := NewProducer(Redis, s.addr, "test")
	c := NewConsumer(Redis, s.addr, "test")
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

func (s *RedisSuite) TestWithMessageChannel() {
	defaultProducerChannel := "test"
	actualChannel := "test1"

	p := NewProducer(Redis, s.addr, defaultProducerChannel)
	c := NewConsumer(Redis, s.addr, actualChannel)
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
		msg.SetChannel(actualChannel)

		s.Require().NoError(p.Produce(context.Background(), msg))

		msg, err := c.Consume(context.Background())
		s.Require().NoError(err)
		s.Require().NotNil(msg)
	}
}

func (s *RedisSuite) TestMessageDecode() {
	p := NewProducer(Redis, s.addr, "test")
	c := NewConsumer(Redis, s.addr, "test")
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

func (s *RedisSuite) TestMessageEncode() {
	p := NewProducer(Redis, s.addr, "test")
	c := NewConsumer(Redis, s.addr, "test")
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

func (s *RedisSuite) setupRedis() (string, *redis.RedisContainer) {
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
