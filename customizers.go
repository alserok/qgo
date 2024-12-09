package qgo

import "time"

// ========================================================
// Kafka

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

func WithOffset(offset int64) Customizer[any] {
	return func(cons any) {
		cons.(*kafkaConsumer).offset = offset
	}
}

func WithPartition(partition int32) Customizer[any] {
	return func(cons any) {
		cons.(*kafkaConsumer).partition = partition
	}
}

func WithRequiredAcks(acks int) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).requiredAcks = acks
	}
}

func WithCompression(comp int) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).compression = comp
	}
}

func WithFlushFrequency(freq time.Duration) Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).flushFreq = freq
	}
}

// =========================================================
// NATS

func WithSubject(subj string) Customizer[any] {
	return func(cons any) {
		switch c := cons.(type) {
		case *natsConsumer:
			c.subject = c.topic + "." + subj
		case *natsPublisher:
			c.subject = c.topic + "." + subj
		default:
			panic("invalid customizer")
		}
	}
}

func WithRetryWait(wait time.Duration) Customizer[any] {
	return func(pub any) {
		pub.(*natsPublisher).retryWait = &wait
	}
}

func WithRetryAttempts(attempts int) Customizer[any] {
	return func(pub any) {
		pub.(*natsPublisher).retryAttempts = &attempts
	}
}
