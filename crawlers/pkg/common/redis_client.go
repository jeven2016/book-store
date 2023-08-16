package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

type Redis struct {
	Client *redis.Client
	config *RedisConfig
}

func NewRedis(redisCfg *RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Address,
		Password:     redisCfg.Password,
		DB:           redisCfg.DefaultDb,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  time.Duration(redisCfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisCfg.WriteTimeout) * time.Second,
		PoolSize:     redisCfg.PoolSize,
		PoolTimeout:  time.Duration(redisCfg.PoolTimeout) * time.Second,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	rd := &Redis{
		Client: client,
		config: redisCfg,
	}

	if err := rd.ensureConsumerGroups(); err != nil {
		return nil, err
	}
	return rd, nil
}

func (rd *Redis) ensureConsumerGroups() error {
	if !rd.config.AutoCreateConsumerGroups {
		return nil
	}

	for streamName, groups := range GetStreamGroupMap() {
		if streamName == "" {
			return nil
		}
		for _, group := range groups {
			if group == "" {
				continue
			}

			//当无法获取到group信息时，创建一个消费group
			if groups, err := rd.Client.XInfoGroups(context.Background(), streamName).Result(); err != nil {
				for _, g := range groups {
					if g.Name == group {
						return nil
					}
				}

				//You can use the XGROUP CREATE command with MKSTREAM option, to create an empty stream
				//XGroupCreate 方法要求先有stream的存在才能创建group
				if err = rd.Client.XGroupCreateMkStream(context.Background(), streamName, group, "0").Err(); err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func (rd *Redis) PublishMessage(ctx context.Context, data interface{}, streamName string) error {
	if data == nil {
		return errors.New(fmt.Sprintf("cannot publis empty data, stream is %v", streamName))
	}
	var json string
	var err error
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.String {
		json = data.(string)
	} else {
		json, err = convertor.ToJson(data)
	}

	if err != nil {
		return errors.New(fmt.Sprintf("unable convert data into json, stream: %s", streamName))
	}

	//just send the json data into stream since it's too complicated to map a struct to map, there would be
	//different kind of exceptions that need to handle
	err = rd.Client.XAdd(ctx, &redis.XAddArgs{
		Stream:     streamName,
		NoMkStream: false,  // * 默认false,当为false时,key不存在，会新建
		MaxLen:     100000, // * 指定stream的最大长度,当队列长度超过上限后，旧消息会被删除，只保留固定长度的新消息
		Approx:     false,  // * 默认false,当为true时,模糊指定stream的长度
		ID:         "*",    // 消息 id，我们使用 * 表示由 redis 生成
		// MinID: "id",            // * 超过阈值，丢弃设置的小于MinID消息id【基本不用】
		// Limit: 1000,            // * 限制长度【基本不用】
		Values: map[string]string{
			RedisStreamDataVar: json,
		}}).Err()

	return err
}

func (rd *Redis) Consume(ctx context.Context, streamName string,
	consumerGroup string, msgChan chan<- string, delErrChan chan<- error) error {
	for {
		if CheckCancel(ctx) {
			return nil
		}
		consumer := streamName + "-consumer-" + strconv.Itoa(int(rand.Uint32()))
		entries, err := rd.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumer,

			Streams: []string{streamName, ">"},
			Count:   2,
			Block:   0,
		}).Result()
		if err != nil {
			return err
		}

		for i := 0; i < len(entries[0].Messages); i++ {
			messageID := entries[0].Messages[i].ID
			jsonData := entries[0].Messages[i].Values[RedisStreamDataVar].(string)
			msgChan <- jsonData
			//必须得ACK，解决redis内存占用高的问题
			rd.Client.XAck(ctx, streamName, consumerGroup, messageID)
			if err = rd.Client.XDel(ctx, streamName, messageID).Err(); err != nil {
				delErrChan <- err
			}
		}
	}
}

// Len returns the current stream length
func (rd *Redis) Len(ctx context.Context, streamName string) (int64, error) {
	streamLen, err := rd.Client.XLen(ctx, streamName).Result()
	if err != nil {
		return 0, err
	}
	return streamLen, err
}
