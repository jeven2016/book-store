package stream

import (
	"context"
	"crawlers/pkg/common"
	"github.com/reugn/go-streams"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
	"github.com/reugn/go-streams/flow"
	"math/rand"
	"strconv"
)

// this file forked and enhanced based on https://github.com/reugn/go-streams since I want to use go-redis/v9 and redis stream feature

//redis stream XADDArg 的values只能是如下格式
// XAddArgs accepts values in the following formats:
//   - XAddArgs.Values = []interface{}{"key1", "value1", "key2", "value2"}
//   - XAddArgs.Values = []string("key1", "value1", "key2", "value2")
//   - XAddArgs.Values = map[string]interface{}{"key1": "value1", "key2": "value2"}

// RedisStreamSource is a Redis Pub/Sub Source
type RedisStreamSource struct {
	ctx           context.Context
	redisClient   *common.Redis
	out           chan interface{}
	streamName    string
	consumerGroup string
}

// NewRedisStreamSource returns a new RedisStreamSource instance
func NewRedisStreamSource(ctx context.Context, client *common.Redis, streamName string,
	consumerGroup string) (*RedisStreamSource, error) {
	var err error
	if err = ensureConsumeGroup(client.Client, consumerGroup, streamName); err != nil {
		return nil, err
	}

	source := &RedisStreamSource{
		ctx:           ctx,
		redisClient:   client,
		out:           make(chan interface{}),
		streamName:    streamName,
		consumerGroup: consumerGroup,
	}

	go source.init(ctx)
	return source, nil
}

func ensureConsumeGroup(client *redis.Client, consumerGroup, stream string) error {
	var groups []redis.XInfoGroup
	var err error

	//当无法获取到group信息时，创建一个消费group
	if groups, err = client.XInfoGroups(context.Background(), stream).Result(); err != nil || len(groups) == 0 {
		var groupExists bool
		for _, g := range groups {
			if g.Name == consumerGroup {
				return nil
			}
		}

		if !groupExists {
			//You can use the XGROUP CREATE command with MKSTREAM option, to create an empty stream
			//XGroupCreate 方法要求先有stream的存在才能创建group
			return client.XGroupCreateMkStream(context.Background(), stream, consumerGroup, "0").Err()
		}
	}

	return nil
}

// init starts the main loop
func (rs *RedisStreamSource) init(ctx context.Context) {
	defer func() {
		close(rs.out)
		if err := rs.redisClient.Client.Close(); err != nil {
			zap.L().Error("failed to close redis connection", zap.Error(err))
		}
		if err := recover(); err != nil {
			zap.S().Errorf("an unexpected error occurs during fetching data form stream, %v", err)
		}
	}()

loop:
	for {

		select {
		case <-ctx.Done():
			break loop
		default:
			rs.fetchFromStream()
		}

	}

}

func (rs *RedisStreamSource) fetchFromStream() {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Errorf("an unexpected error occurs, %v", err)
		}
	}()

	consumer := rs.streamName + ":consumer:" + strconv.Itoa(rand.Int())
	entries, err := rs.redisClient.Client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    rs.consumerGroup,
		Consumer: consumer,

		Streams: []string{rs.streamName, ">"},
		Count:   1,
		Block:   0,
	}).Result()
	if err != nil {
		zap.S().Errorf("failed to handle XReadGroup, %v", err)
		return
	}
	for i := 0; i < len(entries[0].Messages); i++ {
		messageID := entries[0].Messages[i].ID
		jsonData := entries[0].Messages[i].Values[common.RedisStreamDataVar]
		rs.out <- jsonData
		rs.redisClient.Client.XAck(context.Background(), rs.streamName, rs.consumerGroup, messageID)
		zap.S().Info("fetch a message from stream")
	}
}

// Via streams data through the given flow
func (rs *RedisStreamSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(rs, _flow)
	return _flow
}

// Out returns an output channel for sending data
func (rs *RedisStreamSource) Out() <-chan interface{} {
	return rs.out
}

// RedisStreamSink is a Redis Pub/Sub Sink
type RedisStreamSink struct {
	redisClient *common.Redis
	in          chan interface{}
	streamName  string
}

// NewRedisStreamSink returns a new RedisStreamSink instance
func NewRedisStreamSink(ctx context.Context, client *common.Redis, streamName string) *RedisStreamSink {
	sink := &RedisStreamSink{
		client,
		make(chan interface{}),
		streamName,
	}

	go sink.init(ctx)
	return sink
}

// init starts the main loop
func (rs *RedisStreamSink) init(ctx context.Context) {
	defer func() {
		//rs.redisClient.Close()
		if err := recover(); err != nil {
			zap.S().Errorf("an unexpected error occurs, %v", err)
		}
	}()

	for msg := range rs.in {
		if msg == nil {
			continue
		}
		if err := rs.redisClient.PublishMessage(ctx, msg, rs.streamName); err != nil {
			zap.S().Errorf("failed to send a message into stream %v: %v", rs.streamName, err)
		}

	}

}

// In returns an input channel for receiving data
func (rs *RedisStreamSink) In() chan<- interface{} {
	return rs.in
}
