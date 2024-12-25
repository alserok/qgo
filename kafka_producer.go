package qgo

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

func newKafkaProducer(addr, topic string, customs ...Customizer[any]) *kafkaProducer {
	prod := kafkaProducer{
		topic:           topic,
		requiredAcks:    AckWaitForAll,
		compression:     CompressionSnappy,
		returnSuccesses: true,
		returnErrors:    true,
		flushFreq:       time.Millisecond * 500,
	}
	for _, custom := range customs {
		custom(&prod)
	}

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(prod.requiredAcks)
	cfg.Producer.Compression = sarama.CompressionCodec(prod.compression)
	cfg.Producer.Flush.Frequency = prod.flushFreq
	cfg.Producer.Return.Errors = prod.returnErrors
	cfg.Producer.Return.Successes = prod.returnSuccesses
	cfg.Producer.Idempotent = prod.idempotent

	p, err := sarama.NewSyncProducer([]string{addr}, cfg)
	if err != nil {
		panic("failed to init kafkaProducer: " + err.Error())
	}
	prod.p = p

	return &prod
}

type kafkaProducer struct {
	p sarama.SyncProducer

	requiredAcks    int
	compression     int
	returnErrors    bool
	returnSuccesses bool
	idempotent      bool
	flushFreq       time.Duration

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
