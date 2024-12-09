# One Golang library for all message queues

    go get github.com/alserok/qgo

---

* [Kafka](#kafka)
* [Nats](#nats)

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