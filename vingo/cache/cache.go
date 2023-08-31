package cache

import (
	"github.com/lgdzz/vingo-utils/vingo"
	"time"
)

// 从缓存中读取数据
// key 缓存key
// expired 缓存有效期，0则永不过期
// handle 要缓存数据的处理函数
func Fast[T any](key string, expired time.Duration, handle func() T) T {
	return FastRefresh(key, expired, handle, false)
}

func FastRefresh[T any](key string, expired time.Duration, handle func() T, refresh bool) T {
	var result T
	if refresh {
		result = handle()
		Set(key, result, expired)
	} else if dataString := vingo.RedisResult(vingo.Redis.Get(key)); dataString != "" {
		vingo.StringToJson(dataString, &result)
	} else {
		result = handle()
		Set(key, result, expired)
	}
	return result
}

// 设置缓存
func Set(key string, value any, expired time.Duration) {
	_, err := vingo.Redis.Set(key, vingo.JsonToString(value), expired).Result()
	if err != nil {
		panic(err)
	}
}

// 读取缓存
func Get[T any](key string) T {
	var result T
	var dataString = vingo.RedisResult(vingo.Redis.Get(key))
	if dataString != "" {
		vingo.StringToJson(dataString, &result)
	}
	return result
}

// 计算缓存有效期至今日23:59:59
func ExpiredToday() time.Duration {
	now := time.Now()
	expireTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	return time.Duration((expireTime.Unix() - now.Unix() - 1)) * time.Second
}

// 计算缓存有效期至周末23:59:59
func ExpiredWeekEnd() time.Duration {
	now := time.Now()
	// 获取本周的第一天（周一）
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	expireTime := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 23, 59, 59, 0, now.Location()).AddDate(0, 0, 6)
	return time.Duration((expireTime.Unix() - now.Unix() - 1)) * time.Second
}

// 计算缓存有效期至本月最后一天的23:59:59
func ExpiredMomentEnd() time.Duration {
	now := time.Now()
	expireTime := time.Date(now.Year(), now.Month(), 0, 23, 59, 59, 0, now.Location()).AddDate(0, 1, 0)
	return time.Duration((expireTime.Unix() - now.Unix() - 1)) * time.Second
}
