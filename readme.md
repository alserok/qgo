# 📬 One Golang library for all message queues

---

## 🧭 Navigation

* *[Message](#message)*
* *[Kafka](#kafka)*
* *[Nats](#nats)*
* *[Rabbit](#rabbit)*
* *[Redis](#redis)*

---

    go get github.com/alserok/qgo

---

The same interface for all the queues. You may set up everything how you want with
customizers in constructors.

```go
type Producer interface {
	Produce(ctx context.Context, message *Message) error
	Close() error
}

type Consumer interface {
        Consume(ctx context.Context) (*Message, error)
        Close() error
}

```

---

## Message

```go

type Message struct {
        Body      []byte    `json:"body"`
        ID        string    `json:"id"`
        Timestamp time.Time `json:"timestamp"`
	...
}
```

### Encoding/Decoding methods

Methods that allow easily encode and decode `message` body.
```go
// tries to decode message body to target, if fails returns and error
func (m *Message) DecodeBody(target any) error {
    // ...
}

// tries to encode value to message body, if fails returns and error
func (m *Message) EncodeToBody(in any) error {
	// ...
}
```

### Queue methods

Methods that change `message` parameters to dynamically customize their processing.

Kafka
```go
// dynamically change the partition where message will be produced
func (m *Message) SetPartition(part int32) {
    // ...
}
```

Rabbit
```go
// if you set field `autoAcknowledge` to false then you have to call this method 
// after consuming the message
func (m *Message) Ack() error {
	// ...
}

// dynamically set the exchange
func (m *Message) SetExchange(exchange string) {
    // ...
}
```

Nats
```go
// dynamically change the subject where message will be published
func (m *Message) SetSubject(sub string) {
	// ...
}
```

Redis
```go
// dynamically change the channel where message will be published
func (m *Message) SetChannel(ch string) {
    // ...
}
```


---

## Kafka

```go
package main

import (
	"context"
	"fmt"
	"github.com/alserok/qgo"
	"time"
)

func main() {
	cstm := qgo.NewKafkaCustomizer()
	
	p := qgo.NewProducer(qgo.Kafka, "localhost:9092", "test")
	c := qgo.NewConsumer(qgo.Kafka, "localhost:9092", "test", cstm.WithOffset(qgo.OffsetOldest), cstm.WithPartition(0))

	err := p.Produce(context.Background(), &qgo.Message{
		Body:      []byte("body"),
		ID:        "message id",
		Timestamp: time.Now(),
	})
	if err != nil {
		panic("failed to produce message: " + err.Error())
	}

	msg, err := c.Consume(context.Background())
	if err != nil {
		panic("failed to consume message: " + err.Error())
	}
	
	fmt.Println(msg)
}
```

### Customizers

```go
type KafkaCustomizer interface {
    WithOffset(offset int64) Customizer[any]
    WithPartition(partition int32) Customizer[any]
    WithRequiredAcks(acks int) Customizer[any]
    WithCompression(comp int) Customizer[any]
    WithFlushFrequency(freq time.Duration) Customizer[any]
    WithIdempotent() Customizer[any]
}

func NewKafkaCustomizer() KafkaCustomizer {
    // ...
}
```

```go

// sets offset for consumer (by default Newest)
func WithOffset(offset int64) Customizer[any] {
    // ...
}

// sets partition for consumer (by default 0)
func WithPartition(partition int32) Customizer[any] {
    // ...
}

// sets RequiredAcks for producer
func WithRequiredAcks(acks int) Customizer[any]  {
    // ...
}

// sets CompressionCodec for producer
func WithCompression(comp int) Customizer[any]  {
    // ...
}

// sets FlushFrequency for producer
func WithFlushFrequency(freq time.Duration) Customizer[any]  {
    // ...
}

```

---

## Nats

```go
package main

import (
	"context"
	"fmt"
	"github.com/alserok/qgo"
)

func main() {
	cstm := qgo.NewNATSCustomizer()
	
	p := qgo.NewProducer(qgo.Nats, "localhost:4222", "test", cstm.WithSubjects("a"))
	c := qgo.NewConsumer(qgo.Nats, "localhost:4222", "test", cstm.WithSubjects("a"))

	err := p.Produce(context.Background(), &qgo.Message{
		Body:      []byte("body"),
	})
	if err != nil {
		panic("failed to produce message: " + err.Error())
	}

	msg, err := c.Consume(context.Background())
	if err != nil {
		panic("failed to consume message: " + err.Error())
	}
	
	fmt.Println(msg)
}
```

### Customizers

```go
type NatsCustomizer interface {
    WithSubjects(subjects ...string) Customizer[any]
    WithRetryWait(wait time.Duration) Customizer[any]
    WithRetryAttempts(attempts int) Customizer[any]
}

func NewNATSCustomizer() NatsCustomizer {
    // ...
}
```

```go

// consumer
// sets subject to consume (requires only one argument)
//
// producer
// adds stream subjects and sets default subject to produce (the one with index 0)
func WithSubjects(subj ...string) Customizer[any] {
    // ...
}

// sets retry wait for producer
func WithRetryWait(wait time.Duration) Customizer[any] {
    // ...
}

// sets retry attempts for producer
func WithRetryAttempts(attempts int) Customizer[any] {
    // ...
}


```

---

## Rabbit

```go
package main

import (
	"context"
	"fmt"
	"github.com/alserok/qgo"
)

func main() {
	cstm := qgo.NewRabbitCustomizer()
	
	p := qgo.NewProducer(qgo.Rabbit, "localhost:5672", "test")
	c := qgo.NewConsumer(qgo.Rabbit, "localhost:5672", "test", cstm.WithAutoAcknowledgement())

	err := p.Produce(context.Background(), &qgo.Message{
		Body:      []byte("body"),
	})
	if err != nil {
		panic("failed to produce message: " + err.Error())
	}

	msg, err := c.Consume(context.Background())
	if err != nil {
		panic("failed to consume message: " + err.Error())
	}
	
	fmt.Println(msg)
}
```

### Customizers

```go
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
    // ...
}
```

```go

// sets `noWait` field to true
func WithNoWait() Customizer[any] {
    // ...
}

// sets `exclusive` field to true
func WithExclusive() Customizer[any] {
    // ...
}

// sets `autoDelete` field to true
func WithAutoDelete() Customizer[any] {
    // ...
}

// sets `durable` field to true
func WithDurable() Customizer[any] {
    // ...
}

// sets `exchange` field
func WithExchange(exchange string) Customizer[any] {
    // ...
}

// sets `mandatory` field to true
func WithMandatory() Customizer[any] {
    // ...
}

// sets `immediate` field to true
func WithImmediate() Customizer[any] {
    // ...
}

// sets `consumer` field for the rabbitmq Consume method
func WithConsumerTag(tag string) Customizer[any] {
    // ...
}

// sets `immediate` field to true
func WithNoLocal() Customizer[any] {
    // ...
}

// sets `autoAcknowledgement` field to true, by default you have to use Ack method for the message
func WithAutoAcknowledgement() Customizer[any] {
    // ...
}

```

---

## Redis

```go
package main

import (
	"context"
	"fmt"
	"github.com/alserok/qgo"
)

func main() {
	p := qgo.NewProducer(qgo.Redis, "localhost:6379", "test")
	c := qgo.NewConsumer(qgo.Redis, "localhost:6379", "test")

	err := p.Produce(context.Background(), &qgo.Message{
		Body:      []byte("body"),
	})
	if err != nil {
		panic("failed to produce message: " + err.Error())
	}

	msg, err := c.Consume(context.Background())
	if err != nil {
		panic("failed to consume message: " + err.Error())
	}
	
	fmt.Println(msg)
}
```

### Customizers

```go
type RedisCustomizer interface {
	WithPassword(password string) Customizer[any]
	WithDB(db int) Customizer[any]
	WithRetries(retries int) Customizer[any]
	WithTLSConfig(cfg *tls.Config) Customizer[any]
	WithNetwork(net string) Customizer[any]
	WithClientName(name string) Customizer[any]
}

func NewRedisCustomizer() RedisCustomizer {
    // ...
}
```

```go
// sets password for redis client
func WithPassword(password string) Customizer[any] {
	// ...
}

// sets db for redis client
func WithDB(db int) Customizer[any] {
	// ...
}

// sets MaxRetries for redis client options
func (redisCustomizer) WithRetries(retries int) Customizer[any] {
    // ...
}

// sets TLSConfig for redis client options
func (redisCustomizer) WithTLSConfig(cfg *tls.Config) Customizer[any] {
    // ...
}

// sets Network for redis client options
func (redisCustomizer) WithNetwork(net string) Customizer[any] {
    // ...
}

// sets ClientName for redis client options
func (redisCustomizer) WithClientName(name string) Customizer[any] {
    // ...
}
```