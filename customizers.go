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

func WithIdempotent() Customizer[any] {
	return func(prod any) {
		prod.(*kafkaProducer).idempotent = true
	}
}

// =========================================================
// NATS

func WithSubjects(subjects ...string) Customizer[any] {
	return func(cons any) {
		switch cp := cons.(type) {
		case *natsConsumer:
			cp.subject = cp.topic + "." + subjects[0]
		case *natsPublisher:
			for _, subj := range subjects {
				cp.subjects = append(cp.subjects, cp.topic+"."+subj)
			}

			cp.defaultSubject = cp.topic + "." + subjects[0]
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

// =========================================================
// Rabbit

func WithNoWait() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).noWait = true
	}
}

func WithExclusive() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).exclusive = true
	}
}

func WithAutoDelete() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).autoDelete = true
	}
}

func WithDurable() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).durable = true
	}
}

func WithExchange(exchange string) Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).defaultExchange = exchange
	}
}

func WithMandatory() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).mandatory = true
	}
}

func WithImmediate() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).immediate = true
	}
}

func WithConsumerTag(tag string) Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).tag = tag
	}
}

func WithNoLocal() Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).noLocal = true
	}
}

func WithAutoAcknowledgement() Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).autoAcknowledgement = true
	}
}

// =========================================================
// Redis

func WithPassword(password string) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.password = password
		case *redisPublisher:
			ps.password = password
		}
	}
}

func WithDB(db int) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.db = db
		case *redisPublisher:
			ps.db = db
		}
	}
}
