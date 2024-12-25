# One Golang library for all message queues

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

## 🧭 Navigation

* *[Kafka](#kafka)*
* *[Nats](#nats)*
* *[Rabbit](#rabbit)*

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
	p := qgo.NewProducer(qgo.Kafka, "localhost:9092", "test")
	c := qgo.NewConsumer(qgo.Kafka, "localhost:9092", "test", qgo.WithOffset(qgo.OffsetOldest), qgo.WithPartition(0))

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
	p := qgo.NewProducer(qgo.Nats, "localhost:4222", "test", qgo.WithSubject("a"))
	c := qgo.NewConsumer(qgo.Nats, "localhost:4222", "test", qgo.WithSubject("a"))

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

// sets subject for producer and consumer 
func WithSubject(subj string) Customizer[any] {
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
	p := qgo.NewProducer(qgo.Rabbit, "localhost:5672", "test")
	c := qgo.NewConsumer(qgo.Rabbit, "localhost:5672", "test", qgo.WithAutoAcknowledgement())

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