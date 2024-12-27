package qgo

type rabbitCustomizer struct{}

func (rabbitCustomizer) WithNoWait() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).noWait = true
	}
}

func (rabbitCustomizer) WithExclusive() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).exclusive = true
	}
}

func (rabbitCustomizer) WithAutoDelete() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).autoDelete = true
	}
}

func (rabbitCustomizer) WithDurable() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).durable = true
	}
}

func (rabbitCustomizer) WithExchange(exchange string) Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).defaultExchange = exchange
	}
}

func (rabbitCustomizer) WithMandatory() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).mandatory = true
	}
}

func (rabbitCustomizer) WithImmediate() Customizer[any] {
	return func(pub any) {
		pub.(*rabbitPublisher).immediate = true
	}
}

func (rabbitCustomizer) WithConsumerTag(tag string) Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).tag = tag
	}
}

func (rabbitCustomizer) WithNoLocal() Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).noLocal = true
	}
}

func (rabbitCustomizer) WithAutoAcknowledgement() Customizer[any] {
	return func(con any) {
		con.(*rabbitConsumer).autoAcknowledgement = true
	}
}

func (rabbitCustomizer) WithConsumerArguments(args map[string]interface{}) Customizer[any] {
	return func(cons any) {
		cons.(*rabbitConsumer).arguments = args
	}
}
