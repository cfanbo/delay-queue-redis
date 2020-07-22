package queue

type Option func(*Options)

type Options struct {
	topic string
	handler handlerFunc
}

func WithTopic(topic string) Option {
	return func(opts *Options) {
		opts.topic = topic
	}
}

func WithHandler(handler handlerFunc) Option {
	return func(opts *Options) {
		opts.handler = handler
	}
}