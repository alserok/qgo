package qgo

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
)

func newKafkaConsumer(addr, topic string, customs ...Customizer[any]) *kafkaConsumer {
	cfg := sarama.NewConfig()

	cons := kafkaConsumer{
		topic:     topic,
		offset:    sarama.OffsetNewest,
		partition: 0,
	}

	for _, customizer := range customs {
		customizer(&cons)
	}

	c, err := sarama.NewConsumer([]string{addr}, cfg)
	if err != nil {
		panic("failed to init kafka consumer: " + err.Error())
	}
	cons.c = c

	pc, err := c.ConsumePartition(topic, cons.partition, cons.offset)
	if err != nil {
		panic("failed to consume partition '" + topic + "': " + err.Error())
	}
	cons.pc = pc

	return &cons
}

type kafkaConsumer struct {
	c  sarama.Consumer
	pc sarama.PartitionConsumer

	topic     string
	offset    int64
	partition int32
}

func (k *kafkaConsumer) Consume(ctx context.Context) (*Message, error) {
	select {
	case msg, ok := <-k.pc.Messages():
		if !ok {
			return nil, fmt.Errorf("consumer closed")
		}

		return &Message{
			ID:        string(msg.Key),
			Body:      msg.Value,
			Timestamp: msg.Timestamp,
		}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled")
	}
}

func (k *kafkaConsumer) Close() error {
	if err := k.c.Close(); err != nil {
		return fmt.Errorf("failed to close kafka consumer: %w", err)
	}

	if err := k.pc.Close(); err != nil {
		return fmt.Errorf("failed to close kafka partition consumer: %w", err)
	}

	return nil
}
