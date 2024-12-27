package qgo

import "time"

type natsCustomizer struct{}

func (natsCustomizer) WithSubjects(subjects ...string) Customizer[any] {
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

func (natsCustomizer) WithRetryWait(wait time.Duration) Customizer[any] {
	return func(pub any) {
		pub.(*natsPublisher).retryWait = &wait
	}
}

func (natsCustomizer) WithRetryAttempts(attempts int) Customizer[any] {
	return func(pub any) {
		pub.(*natsPublisher).retryAttempts = &attempts
	}
}
