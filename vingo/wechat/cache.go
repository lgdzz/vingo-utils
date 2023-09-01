package wechat

import (
	"github.com/lgdzz/vingo-utils/vingo/db/redis"
	"time"
)

type Cache struct {
}

// Get 获取一个值
func (s *Cache) Get(key string) interface{} {
	result, err := redis.Client.Get(key).Result()
	if err != nil {
		return nil
	}
	return result
}

// Set 设置一个值
func (s *Cache) Set(key string, val interface{}, timeout time.Duration) error {
	return redis.Client.Set(key, val, timeout).Err()
}

// IsExist 判断key是否存在
func (s *Cache) IsExist(key string) bool {
	result, _ := redis.Client.Exists(key).Result()
	return result > 0
}

// Delete 删除
func (s *Cache) Delete(key string) error {
	return redis.Client.Del(key).Err()
}

var CacheApi = &Cache{}
