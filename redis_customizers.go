package qgo

import "crypto/tls"

type redisCustomizer struct{}

func (redisCustomizer) WithPassword(password string) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.password = password
		case *redisPublisher:
			ps.password = password
		}
	}
}

func (redisCustomizer) WithDB(db int) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.db = db
		case *redisPublisher:
			ps.db = db
		}
	}
}

func (redisCustomizer) WithRetries(retries int) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.maxRetries = retries
		case *redisPublisher:
			ps.maxRetries = retries
		}
	}
}

func (redisCustomizer) WithTLSConfig(cfg *tls.Config) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.tlsConfig = cfg
		case *redisPublisher:
			ps.tlsConfig = cfg
		}
	}
}

func (redisCustomizer) WithNetwork(net string) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.network = net
		case *redisPublisher:
			ps.network = net
		}
	}
}

func (redisCustomizer) WithClientName(name string) Customizer[any] {
	return func(pubSub any) {
		switch ps := pubSub.(type) {
		case *redisConsumer:
			ps.clientName = name
		case *redisPublisher:
			ps.clientName = name
		}
	}
}
