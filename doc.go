//package main
//
//import (
//	"context"
//	"fmt"
//	"log"
//	"time"
//
//	queue "github.com/cfanbo/delay-queue-redis"
//	"github.com/go-redis/redis/v8"
//)
//
//var redisClient *redis.Client
//
//type Msg struct {
//	MsgId   int    `json:"msg_id"`
//	MsgBody string `json:"body"`
//	UserId  int    `json:"uid"`
//}
//
//func handerFunc(msg queue.Message) {
//	fmt.Println("消费一条消息：=========")
//	fmt.Printf("%#v\n", msg)
//
//	// 转map
//	m := msg.Body.(map[string]interface{})
//	fmt.Println(m["msg_id"], m["body"], m["uid"])
//}
//
//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	redisClient = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//
//	_, err := redisClient.Ping(ctx).Result()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 创建延时队列
//	q := queue.NewQueue(ctx, redisClient, queue.WithTopic("test-topic"), queue.WithHandler(handerFunc))
//	q.Start()
//
//	// 创建消息实体对象
//	ticker := time.NewTicker(time.Second * 1)
//	go func(ticker *time.Ticker) {
//		defer ticker.Stop()
//
//		for {
//			select {
//			case <-ticker.C:
//				message := Msg{100, "abc", 43}
//				msg := queue.NewMessage("", time.Now().Add(time.Second*8), message)
//
//				// 发布
//				_, err = q.Publish(msg)
//				if err != nil {
//					log.Fatal(err)
//				}
//				fmt.Println("发布成功一条消息")
//			}
//		}
//
//	}(ticker)
//
//	// 手动延时10秒后退出
//	time.Sleep(time.Second * 10)
//	cancel()
//}

package queue
