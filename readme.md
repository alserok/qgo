# One Golang library for all message queues

    go get github.com/alserok/qgo

---

* [Kafka](#kafka)

---

## Kafka producer and consumer

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
	c := qgo.NewConsumer(qgo.Kafka, "localhost:9092", "test", qgo.KafkaConsumerConfig{
		Partition: 0,
		Offset:    qgo.OffsetNewest,
	})

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

```go
type KafkaConsumerConfig struct {
	Partition int32 // provide your partition or will be chosen default
	Offset    int64 // provide your offset or will be chosen `newest`
}
```