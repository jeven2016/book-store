package stream

//
//import (
//	"context"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/reugn/go-streams/flow"
//)
//
//const sourceStream = "sourceStream"
//const testConsumerGroup = "testConsumerGroup"
//
//const destStream = "destStream"
//
//// docker exec -it pubsub bash
//// https://redis.io/topics/pubsub
//func TestStream(t *testing.T) {
//	ctx, cancelFunc := context.WithCancel(context.Background())
//
//	timer := time.NewTimer(2 * time.Minute)
//	go func() {
//		<-timer.C
//		cancelFunc()
//	}()
//
//	source, err := NewRedisStreamSource(ctx, redisClient(), sourceStream, testConsumerGroup)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	toUpperMapFlow := flow.NewMap(toUpper, 1)
//
//	sink := NewRedisStreamSink(ctx, redisClient(), destStream)
//
//	source.
//		Via(toUpperMapFlow).
//		To(sink)
//
//}
//
//var toUpper = func(msg map[string]interface{}) map[string]interface{} {
//	for k, v := range msg {
//		msg[k] = strings.ToUpper(v.(string))
//	}
//	log.Printf("map msg: %v", msg)
//	return msg
//}
//
//func redisClient() (client *redis.Client) {
//
//	client = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "pwd", // no password set
//		DB:       0,     // use default DB
//		PoolSize: 3,     // 默认一个CPU 10个连接
//	})
//	_, err := client.Ping(context.Background()).Result()
//	if err != nil {
//		log.Fatal("Unable to connect to Redis", err)
//	}
//	return client
//}
