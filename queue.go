package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

const (
	HashSuffix = ":hash"
	SetSuffix = ":set"
)

var once sync.Once

type Queue struct {
	// ctx
	ctx context.Context

	// redis
	redis *redis.Client
	topic string

	// producter and consumer
	producer *producer
	consumer *consumer
}

func NewQueue(ctx context.Context, redis *redis.Client, opts ...Option) *Queue{
	var queue *Queue

	once.Do(func() {
		defaultOptions := Options{
			topic: "topic",
			handler: defaultHander,
		}

		for _, apply := range opts {
			apply(&defaultOptions)
		}

		queue = &Queue{
			ctx: ctx,
			redis:    redis,
			topic:    defaultOptions.topic,
			producer: NewProducer(ctx),
			consumer: NewConsumer(ctx, defaultOptions.handler),
		}
	})

	return queue
}

func (q *Queue)Start() {
	go q.consumer.listen(q.redis, q.topic)
}

func (q *Queue)Publish(msg *Message) (int64, error) {
	return q.producer.publish(q.redis, q.topic, msg)
}
