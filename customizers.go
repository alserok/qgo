package qgo

import "time"

const (
	OffsetNewest = -1
	OffsetOldest = -2
)

// ========================================================
// Kafka

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
