package qgo

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
)

func newKafkaConsumer(addr, topic string, consumerCfg ...any) *kafkaConsumer {
	cfg := sarama.NewConfig()

	c, err := sarama.NewConsumer([]string{addr}, cfg)
	if err != nil {
		panic("failed to init kafka consumer: " + err.Error())
	}

	var (
		partition int32
		offset    = sarama.OffsetNewest
	)
	if len(consumerCfg) == 1 {
		consCfg := consumerCfg[0].(KafkaConsumerConfig)

		if consCfg.Partition != 0 {
			partition = consCfg.Partition
		}

		if consCfg.Offset == OffsetOldest {
			offset = consCfg.Offset
		}
	}

	pc, err := c.ConsumePartition(topic, partition, offset)
	if err != nil {
		panic("failed to consume partition '" + topic + "': " + err.Error())
	}

	return &kafkaConsumer{
		c:         c,
		pc:        pc,
		topic:     topic,
		offset:    0,
		partition: 0,
	}
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
