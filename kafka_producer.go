package qgo

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

func newKafkaProducer(addr, topic string) *kafkaProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer([]string{addr}, cfg)
	if err != nil {
		panic("failed to init kafkaProducer: " + err.Error())
	}

	return &kafkaProducer{
		p:     p,
		topic: topic,
	}
}

type kafkaProducer struct {
	p sarama.SyncProducer

	topic string
}

func (p *kafkaProducer) Produce(ctx context.Context, message *Message) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context canceled: %w", ctx.Err())
	}

	_, _, err := p.p.SendMessage(&sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(message.ID),
		Value:     sarama.StringEncoder(message.Body),
		Timestamp: message.Timestamp,
		Partition: message.Partition,
	})
	if err != nil {
		return fmt.Errorf("kafka kafkaProducer failed to send message: %w", err)
	}

	return nil
}

func (p *kafkaProducer) Close() error {
	if err := p.p.Close(); err != nil {
		return fmt.Errorf("failed to close kafka kafkaProducer: %w", err)
	}

	return nil
}
