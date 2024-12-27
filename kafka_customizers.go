package qgo

import "time"

const (
	CompressionNone = iota
	CompressionGZIP
	CompressionSnappy
	CompressionLZ4
	CompressionZSTD

	OffsetNewest = -1
	OffsetOldest = -2

	AckWaitForAll   = -1
	AckWaitForLocal = 1
	AckNoResponse   = 0
)

type kafkaCustomizer struct{}

func (kafkaCustomizer) WithOffset(offset int64) Customizer[any] {
	return func(cons any) {
		cons.(*kafkaConsumer).offset = offset
	}
}

func (kafkaCustomizer) WithPartition(partition int32) Customizer[any] {
	return func(cons any) {
		cons.(*kafkaConsumer).partition = partition
	}
}

func (kafkaCustomizer) WithRequiredAcks(acks int) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).requiredAcks = acks
	}
}

func (kafkaCustomizer) WithCompression(comp int) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).compression = comp
	}
}

func (kafkaCustomizer) WithFlushFrequency(freq time.Duration) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).flushFreq = freq
	}
}

func (kafkaCustomizer) WithIdempotent() Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).idempotent = true
	}
}
