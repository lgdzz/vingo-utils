package redis

import (
	"github.com/go-redis/redis"
	"github.com/lgdzz/vingo-utils/vingo"
	"time"
)

func Get[T any](key string) *T {
	text, err := Client.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	} else {
		var data T
		vingo.StringToJson(text, &data)
		return &data
	}
}

func Set(key string, value any, expiration time.Duration) string {
	result, err := Client.Set(key, vingo.JsonToString(value), expiration).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func HSet(key string, field string, value any) bool {
	result, err := Client.HSet(key, field, vingo.JsonToString(value)).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func HGet[T any](key string, field string) *T {
	text, err := Client.HGet(key, field).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	}
	var data T
	vingo.StringToJson(text, &data)
	return &data
}

func Del(key ...string) int64 {
	result, err := Client.Del(key...).Result()
	if err != nil {
		panic(err)
	}
	return result
}
