package common

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
)

type valueProvider func() (string, error)

// GetKey get a value from cache by key if presents otherwise set by value provider
func GetKey(ctx context.Context, cache *redis.Client,
	key string, callback valueProvider) (string, error) {
	value, err := cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			value, err = callback()
			if err != nil {
				return "", err
			}

			if _, err = cache.Set(ctx, key, value, GenExpireTime()).Result(); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return value, err
}

func Exists(ctx context.Context, key string, searchMongoFunc func() (any, error)) (bool, error) {
	sys := GetSystem()
	rd := sys.RedisClient.Client

	//check if it exists in redis
	if result, err := rd.Exists(ctx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	} else if result > 0 {
		return true, nil
	}

	//check if it exists in mongo
	if val, err := searchMongoFunc(); err != nil {
		return false, err
	} else {
		var realExists bool
		//判断值
		valType := reflect.TypeOf(val)
		if valType.Kind() == reflect.Ptr && valType.Elem() != nil {
			realExists = true
		} else if valType.Kind() == reflect.Bool {
			realExists = reflect.ValueOf(val).Bool()
		}

		if realExists {
			//cache in redis if it exists in mongo
			if _, err = rd.Set(ctx, key, "1", GenExpireTime()).Result(); err != nil {
				return false, err
			}
		}
		return realExists, nil
	}
}

func GenKey(keys ...string) string {
	return strings.Join(keys, ":")
}
