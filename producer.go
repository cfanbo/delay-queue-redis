package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type producer struct {
	ctx context.Context
}

func NewProducer(ctx context.Context) *producer {
	return &producer{
		ctx: ctx,
	}
}

// Publish
func (p *producer)publish(redisClient *redis.Client, topic string, msg *Message) (int64, error) {
	z := &redis.Z{
		Score: msg.GetScore(),
		Member: msg.GetId(),
	}

	// stored sets 写入
	key := topic + SetSuffix
	n, err := redisClient.ZAdd(p.ctx, key, z).Result()
	if err != nil {
		return n, err
	}

	// hashes 写入
	key = topic + HashSuffix
	return redisClient.HSet(p.ctx, key, msg.GetId(), msg).Result()
}