package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

type handlerFunc func(msg Message)

func defaultHander(msg Message) {
	fmt.Println(msg)
}

type consumer struct {
	ctx context.Context
	duration time.Duration
	ch chan []string
	handler handlerFunc
}

func NewConsumer(ctx context.Context, handler handlerFunc) *consumer {
	return &consumer{
		ctx: ctx,
		duration: time.Second,
		ch: make(chan []string, 1000),
		handler: handler,
	}
}

func (c *consumer)listen(redisClient *redis.Client, topic string) {
	// 从 Hashes 中获取数据并处理
	go func() {
		for {
			select {
			case ret := <-c.ch:
				// 批量从hashes中获取数据信息
				key := topic + HashSuffix
				result, err := redisClient.HMGet(c.ctx, key, ret...).Result()
				if err != nil {
					log.Println(err)
				}

				if len(result) > 0 {
					redisClient.HDel(c.ctx, key, ret...)
				}

				msg := Message{}
				for _, v:=range result {
					// 由于hashes 和 scoreSet 非事务操作，会出现删除了set但hashes未删除的情况
					if v == nil {
						continue
					}
					str := v.(string)
					json.Unmarshal([]byte(str), &msg)

					// 处理逻辑
					go c.handler(msg)
				}

			}
		}
	}()

	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			log.Println("consumer quit:", c.ctx.Err())
			return
		case <-ticker.C:
			// read data from redis
			min := strconv.Itoa(0)
			max := strconv.Itoa(int(time.Now().Unix()))
			opt := &redis.ZRangeBy{
				Min: min,
				Max: max,
			}

			key := topic + SetSuffix
			result, err := redisClient.ZRangeByScore(c.ctx, key, opt).Result()
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println(result)

			// 获取到数据
			if len(result) > 0 {
				// 从 sorted sets 中移除数据
				redisClient.ZRemRangeByScore(c.ctx, key, min, max)

				// 写入 chan, 进行hashes处理
				c.ch <- result
			}
		}
	}
}
